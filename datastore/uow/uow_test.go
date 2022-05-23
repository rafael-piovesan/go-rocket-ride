//go:build integration
// +build integration

package uow

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/audit"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tabbed/pqtype"
)

func TestSQLStore(t *testing.T) {
	ctx := context.Background()

	// database up
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer func() { _ = terminate(ctx) }()

	// migrations up
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	// test fixtures up
	userID := int64(gofakeit.Number(0, 1000))
	keyID := int64(gofakeit.Number(0, 1000))

	err = testfixtures.Load(
		dsn,
		[]string{
			"db/fixtures/users",
			"db/fixtures/idempotency_keys",
		},
		map[string]interface{}{
			"UserId":    userID,
			"UserEmail": gofakeit.Email(),
			"KeyId":     keyID,
		},
	)
	require.NoError(t, err)

	// connect to database
	db, _ := db.Connect(config.Config{DBSource: dsn})
	store := New(db)
	rides := datastore.NewRide(db)

	// test entities
	ride := &entity.Ride{
		IdempotencyKeyID: &keyID,
		OriginLat:        0.0,
		OriginLon:        0.0,
		TargetLat:        0.0,
		TargetLon:        0.0,
		UserID:           userID,
	}

	ip := pqtype.CIDR{}
	err = ip.Scan(gofakeit.IPv4Address())
	require.NoError(t, err)

	ar := &entity.AuditRecord{
		Action:       audit.ActionCreateRide,
		Data:         []byte("{\"data\": \"foo\"}"),
		OriginIP:     ip.IPNet.String(),
		ResourceID:   int64(gofakeit.Number(0, 1000)),
		ResourceType: audit.ResourceTypeRide,
		UserID:       userID,
	}

	t.Run("Rollback on error", func(t *testing.T) {
		_, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
		require.ErrorIs(t, err, data.ErrRecordNotFound)

		err = store.Do(ctx, func(ds UnitOfWorkStore) error {
			err := ds.Rides().Save(ctx, ride)
			require.NoError(t, err)

			err = ds.AuditRecords().Save(ctx, ar)
			require.NoError(t, err)

			return errors.New("error rollback")
		})

		if assert.EqualError(t, err, "error rollback") {
			_, err = rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
			assert.ErrorIs(t, err, data.ErrRecordNotFound)
		}
	})

	t.Run("Rollback on panic with error", func(t *testing.T) {
		_, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
		require.ErrorIs(t, err, data.ErrRecordNotFound)

		defer func() {
			p := recover()
			if assert.NotNil(t, p) {
				_, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
				assert.ErrorIs(t, err, data.ErrRecordNotFound)
			}
		}()

		_ = store.Do(ctx, func(ds UnitOfWorkStore) error {
			err := ds.Rides().Save(ctx, ride)
			require.NoError(t, err)

			err = ds.AuditRecords().Save(ctx, ar)
			require.NoError(t, err)

			panic(errors.New("panic rollback"))
		})
	})

	t.Run("Rollback on panic without error", func(t *testing.T) {
		_, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
		require.ErrorIs(t, err, data.ErrRecordNotFound)

		defer func() {
			p := recover()
			if assert.NotNil(t, p) && assert.Equal(t, "panic rollback", p) {
				_, err = rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
				assert.ErrorIs(t, err, data.ErrRecordNotFound)
			}
		}()

		err = store.Do(ctx, func(ds UnitOfWorkStore) error {
			err := ds.Rides().Save(ctx, ride)
			require.NoError(t, err)

			err = ds.AuditRecords().Save(ctx, ar)
			require.NoError(t, err)

			panic("panic rollback")
		})
	})

	t.Run("Rollback on context canceled", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		_, err := rides.FindOne(cancelCtx, datastore.RideWithIdemKeyID(keyID))
		require.ErrorIs(t, err, data.ErrRecordNotFound)

		err = store.Do(ctx, func(ds UnitOfWorkStore) error {
			err := ds.Rides().Save(cancelCtx, ride)
			require.NoError(t, err)

			cancel()

			// this call should return an error due to the canceled ctx
			err = ds.AuditRecords().Save(cancelCtx, ar)
			return err
		})

		if assert.EqualError(t, err, "context canceled") {
			_, err = rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
			assert.ErrorIs(t, err, data.ErrRecordNotFound)
		}
	})

	t.Run("Commit on success", func(t *testing.T) {
		_, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
		require.ErrorIs(t, err, data.ErrRecordNotFound)

		err = store.Do(ctx, func(ds UnitOfWorkStore) error {
			err := ds.Rides().Save(ctx, ride)
			require.NoError(t, err)

			err = ds.AuditRecords().Save(ctx, ar)
			require.NoError(t, err)

			return nil
		})

		if assert.NoError(t, err) {
			res, err := rides.FindOne(ctx, datastore.RideWithIdemKeyID(keyID))
			if assert.NoError(t, err) {
				ride.ID = res.ID
				assert.Equal(t, *ride, res)
			}
		}
	})
}
