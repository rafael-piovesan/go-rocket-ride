//go:build integration
// +build integration

package datastore

import (
	"context"
	"testing"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/stagedjob"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStagedJob(t *testing.T) {
	ctx := context.Background()

	// database up
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer func() { _ = terminate(ctx) }()

	// migrations up
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	// connect to database
	db, _ := db.Connect(config.Config{DBSource: dsn})
	store := NewStagedJob(db)

	t.Run("Create Staged Job", func(t *testing.T) {
		sj := &entity.StagedJob{
			JobName: stagedjob.JobNameSendReceipt,
			JobArgs: []byte("{\"data\": \"foo\"}"),
		}

		err := store.Save(ctx, sj)

		if assert.NoError(t, err) {
			assert.Greater(t, sj.ID, int64(0))
		}
	})
}
