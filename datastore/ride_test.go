//go:build integration
// +build integration

package datastore

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRide(t *testing.T) {
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

	// conntect to database
	db, _ := db.Connect(config.Config{DBSource: dsn})
	store := NewRide(db)

	// test entity
	ride := &entity.Ride{
		IdempotencyKeyID: &keyID,
		OriginLat:        0.0,
		OriginLon:        0.0,
		TargetLat:        0.0,
		TargetLon:        0.0,
		UserID:           userID,
	}

	t.Run("Ride not found", func(t *testing.T) {
		_, err := store.FindOne(ctx, RideWithIdemKeyID(keyID))
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("Create Ride", func(t *testing.T) {
		err := store.Save(ctx, ride)
		if assert.NoError(t, err) {
			assert.Greater(t, ride.ID, int64(0))
		}
	})

	t.Run("Update Ride", func(t *testing.T) {
		stripeID := gofakeit.UUID()
		ride.StripeChargeID = &stripeID

		err := store.Update(ctx, ride)
		assert.NoError(t, err)
	})

	t.Run("Get Ride By Idempotency Key ID", func(t *testing.T) {
		res, err := store.FindOne(ctx, RideWithIdemKeyID(keyID))
		if assert.NoError(t, err) {
			assert.Equal(t, *ride, res)
		}
	})
}
