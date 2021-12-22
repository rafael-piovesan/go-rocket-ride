//go:build integration
// +build integration

package datastore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestRide(t *testing.T) {
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

	// conntect to database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	store := NewStore(db)

	// test entity
	ride := &entity.Ride{
		IdempotencyKeyID: &keyID,
		OriginLat:        "0.0000000000",
		OriginLon:        "0.0000000000",
		TargetLat:        "0.0000000000",
		TargetLon:        "0.0000000000",
		UserID:           userID,
	}

	t.Run("Ride not found", func(t *testing.T) {
		_, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Create Ride", func(t *testing.T) {
		res, err := store.CreateRide(ctx, ride)
		if assert.NoError(t, err) {
			ride.ID = res.ID
			assert.Equal(t, ride, res)
		}
	})

	t.Run("Update Ride", func(t *testing.T) {
		stripeID := gofakeit.UUID()
		ride.StripeChargeID = &stripeID

		res, err := store.UpdateRide(ctx, ride)
		if assert.NoError(t, err) {
			assert.Equal(t, ride, res)
		}
	})

	t.Run("Get Ride By Idempotency Key ID", func(t *testing.T) {
		res, err := store.GetRideByIdempotencyKeyID(ctx, keyID)
		if assert.NoError(t, err) {
			assert.Equal(t, ride, res)
		}
	})
}
