//go:build integration
// +build integration

package bun

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdempotencyKey(t *testing.T) {
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
	idemKey := gofakeit.UUID()
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": gofakeit.Email(),
	})
	require.NoError(t, err)

	// conntect to database
	store, err := NewStore(dsn)
	require.NoError(t, err)

	// test entity
	now := time.Now()
	ttime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	ik := &entity.IdempotencyKey{
		IdempotencyKey: idemKey,
		LastRunAt:      ttime,
		LockedAt:       &ttime,
		RequestMethod:  gofakeit.HTTPMethod(),
		RequestParams:  []byte("{\"data\": \"foo\"}"),
		RequestPath:    fmt.Sprintf("/%s/%s", gofakeit.AnimalType(), gofakeit.Animal()),
		RecoveryPoint:  idempotency.RecoveryPointStarted,
		UserID:         userID,
	}

	t.Run("Idempotency Key not found", func(t *testing.T) {
		_, err := store.GetIdempotencyKey(ctx, idemKey, userID)
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})

	t.Run("Create Idempotency Key", func(t *testing.T) {
		res, err := store.CreateIdempotencyKey(ctx, ik)
		if assert.NoError(t, err) {
			ik.ID = res.ID
			assert.Equal(t, ik, res)
		}
	})

	t.Run("Update Idempotency Key", func(t *testing.T) {
		now = time.Now()

		rps := []idempotency.RecoveryPoint{
			idempotency.RecoveryPointCreated,
			idempotency.RecoveryPointCharged,
			idempotency.RecoveryPointFinished,
		}

		idx := gofakeit.Number(0, len(rps)-1)

		resCode := idempotency.ResponseCodeOK
		resBody := idempotency.ResponseBody{Message: "OK"}

		ik.LastRunAt = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
		ik.LockedAt = nil
		ik.RecoveryPoint = rps[idx]
		ik.ResponseCode = &resCode
		ik.ResponseBody = &resBody

		res, err := store.UpdateIdempotencyKey(ctx, ik)
		if assert.NoError(t, err) {
			assert.Equal(t, ik, res)
		}
	})

	t.Run("Get Idempotency Key", func(t *testing.T) {
		res, err := store.GetIdempotencyKey(ctx, idemKey, userID)
		if assert.NoError(t, err) {
			assert.Equal(t, ik, res)
		}
	})
}
