//go:build integration
// +build integration

package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdempotencyKey(t *testing.T) {
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
	idemKey := gofakeit.UUID()
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": gofakeit.Email(),
	})
	require.NoError(t, err)

	// conntect to database
	db, _ := db.Connect(config.Config{DBSource: dsn})
	store := NewIdempotencyKey(db)

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
		_, err := store.FindOne(ctx, IdemKeyWithKey(idemKey), IdemKeyWithUserID(userID))
		assert.ErrorIs(t, err, repo.ErrRecordNotFound)
	})

	t.Run("Create Idempotency Key", func(t *testing.T) {
		err := store.Save(ctx, ik)
		if assert.NoError(t, err) {
			assert.Greater(t, ik.ID, int64(0))
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

		err := store.Update(ctx, ik)
		assert.NoError(t, err)
	})

	t.Run("Get Idempotency Key", func(t *testing.T) {
		res, err := store.FindOne(ctx, IdemKeyWithKey(idemKey), IdemKeyWithUserID(userID))
		if assert.NoError(t, err) {
			assert.Equal(t, *ik, res)
		}
	})
}
