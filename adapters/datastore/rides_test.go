//go:build integration
// +build integration

package datastore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/require"
)

func TestRides(t *testing.T) {
	ctx := context.Background()

	// create a new Postgres container to be used in tests
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer terminate(ctx)

	// run all migrations
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	userID := int64(gofakeit.Number(0, 10000))
	keyID := int64(gofakeit.Number(0, 10000))
	rideID := int64(gofakeit.Number(0, 10000))

	// load database test fixtures separately
	err = testfixtures.Load(
		dsn,
		[]string{
			"db/fixtures/users",
			"db/fixtures/idempotency_keys",
			"db/fixtures/rides",
		},
		map[string]interface{}{
			"UserId":    userID,
			"UserEmail": gofakeit.Email(),
			"KeyId":     keyID,
			"RideId":    rideID,
		},
	)
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	testQueries := New(db)

	t.Run("Get Ride by ID", func(t *testing.T) {
		r, err := testQueries.GetRideByID(ctx, rideID)
		require.NoError(t, err)

		require.Equal(t, rideID, r.ID)
		require.Equal(t, userID, r.UserID)
	})

	t.Run("Get Ride by Idempotency Key ID", func(t *testing.T) {
		k := sql.NullInt64{Int64: keyID, Valid: true}
		r, err := testQueries.GetRideByIdempotencyKeyID(ctx, k)
		require.NoError(t, err)

		require.Equal(t, rideID, r.ID)
		require.Equal(t, userID, r.UserID)
	})

	t.Run("Update Ride", func(t *testing.T) {
		ride := UpdateRideParams{
			ID:             rideID,
			StripeChargeID: sql.NullString{String: gofakeit.UUID(), Valid: true},
		}
		r, err := testQueries.UpdateRide(ctx, ride)
		require.NoError(t, err)

		require.Equal(t, ride.ID, r.ID)
		require.Equal(t, ride.StripeChargeID, r.StripeChargeID)
	})
}
