//go:build integration
// +build integration

package bun

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/audit"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tabbed/pqtype"
)

func TestSQLStore(t *testing.T) {
	ctx := context.Background()

	// database up
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer terminate(ctx)

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
	store, err := NewStore(dsn)
	require.NoError(t, err)

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
		_, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		require.ErrorIs(t, err, entity.ErrNotFound)

		err = store.Atomic(ctx, func(ds rocketride.Datastore) error {
			_, err := ds.CreateRide(ctx, ride)
			require.NoError(t, err)

			_, err = ds.CreateAuditRecord(ctx, ar)
			require.NoError(t, err)

			return errors.New("error rollback")
		})

		if assert.EqualError(t, err, "error rollback") {
			_, err = store.GetRideByIdempotencyKeyID(ctx, keyID)
			assert.ErrorIs(t, err, entity.ErrNotFound)
		}
	})

	t.Run("Rollback on panic with error", func(t *testing.T) {
		_, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		require.ErrorIs(t, err, entity.ErrNotFound)

		err = store.Atomic(ctx, func(ds rocketride.Datastore) error {
			_, err := ds.CreateRide(ctx, ride)
			require.NoError(t, err)

			_, err = ds.CreateAuditRecord(ctx, ar)
			require.NoError(t, err)

			panic(errors.New("panic rollback"))
		})

		if assert.EqualError(t, err, "panic err: panic rollback") {
			_, err = store.GetRideByIdempotencyKeyID(ctx, keyID)
			assert.ErrorIs(t, err, entity.ErrNotFound)
		}
	})

	t.Run("Rollback on panic without error", func(t *testing.T) {
		_, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		require.ErrorIs(t, err, entity.ErrNotFound)

		defer func() {
			p := recover()
			if assert.NotNil(t, p) && assert.Equal(t, "panic rollback", p) {
				_, err = store.GetRideByIdempotencyKeyID(ctx, keyID)
				assert.ErrorIs(t, err, entity.ErrNotFound)
			}
		}()

		err = store.Atomic(ctx, func(ds rocketride.Datastore) error {
			_, err := ds.CreateRide(ctx, ride)
			require.NoError(t, err)

			_, err = ds.CreateAuditRecord(ctx, ar)
			require.NoError(t, err)

			panic("panic rollback")
		})
	})

	t.Run("Rollback on context canceled", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		_, err := store.GetRideByIdempotencyKeyID(cancelCtx, keyID)
		require.ErrorIs(t, err, entity.ErrNotFound)

		err = store.Atomic(ctx, func(ds rocketride.Datastore) error {
			_, err := ds.CreateRide(cancelCtx, ride)
			require.NoError(t, err)

			cancel()

			// this call should return an error due to the canceled ctx
			_, err = ds.CreateAuditRecord(cancelCtx, ar)
			return err
		})

		if assert.EqualError(t, err, "context canceled") {
			_, err = store.GetRideByIdempotencyKeyID(ctx, keyID)
			assert.ErrorIs(t, err, entity.ErrNotFound)
		}
	})

	t.Run("Commit on success", func(t *testing.T) {
		_, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		require.ErrorIs(t, err, entity.ErrNotFound)

		err = store.Atomic(ctx, func(ds rocketride.Datastore) error {
			_, err := ds.CreateRide(ctx, ride)
			require.NoError(t, err)

			_, err = ds.CreateAuditRecord(ctx, ar)
			require.NoError(t, err)

			return nil
		})

		if assert.NoError(t, err) {
			res, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
			if assert.NoError(t, err) {
				ride.ID = res.ID
				assert.Equal(t, ride, res)
			}
		}
	})

}
