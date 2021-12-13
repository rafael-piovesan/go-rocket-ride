//go:build integration
// +build integration

package datastore

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/require"
	"github.com/tabbed/pqtype"
)

func TestIdempotencyKeys(t *testing.T) {
	ctx := context.Background()

	// create a new Postgres container to be used in tests
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer terminate(ctx)

	// run all migrations
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	userID := int64(gofakeit.Number(0, 10000))

	// load database test fixtures separately
	err = testfixtures.Load(
		dsn,
		[]string{
			"db/fixtures/users",
		},
		map[string]interface{}{
			"UserId":    userID,
			"UserEmail": gofakeit.Email(),
		},
	)
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	testQueries := New(db)

	var ik IdempotencyKey

	t.Run("Create Idempotency Key", func(t *testing.T) {
		i := CreateIdempotencyKeyParams{
			CreatedAt:      time.Now(),
			IdempotencyKey: gofakeit.UUID(),
			LastRunAt:      time.Now(),
			LockedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			RequestMethod:  http.MethodPost,
			RequestParams:  []byte("{}"),
			RequestPath:    "/v1/api",
			RecoveryPoint:  "START",
			UserID:         userID,
		}

		r, err := testQueries.CreateIdempotencyKey(ctx, i)
		require.NoError(t, err)

		require.Equal(t, i.IdempotencyKey, r.IdempotencyKey)
		require.Equal(t, i.RequestMethod, r.RequestMethod)
		require.Equal(t, i.RequestParams, r.RequestParams)
		require.Equal(t, i.RequestPath, r.RequestPath)
		require.Equal(t, i.RecoveryPoint, r.RecoveryPoint)
		require.Equal(t, i.UserID, r.UserID)

		ik = r
	})

	t.Run("Update Idempotency Key", func(t *testing.T) {
		i := UpdateIdempotencyKeyParams{
			ID:            ik.ID,
			LastRunAt:     time.Now(),
			LockedAt:      sql.NullTime{},
			ResponseCode:  sql.NullInt32{Int32: http.StatusOK, Valid: true},
			ResponseBody:  pqtype.NullRawMessage{RawMessage: []byte("{}"), Valid: true},
			RecoveryPoint: "FINISH",
		}
		r, err := testQueries.UpdateIdempotencyKey(ctx, i)
		require.NoError(t, err)

		require.Equal(t, i.LockedAt, r.LockedAt)
		require.Equal(t, i.ResponseCode, r.ResponseCode)
		require.Equal(t, i.ResponseBody, r.ResponseBody)
		require.Equal(t, i.RecoveryPoint, r.RecoveryPoint)

	})
}
