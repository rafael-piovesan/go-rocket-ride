//go:build integration
// +build integration

package datastore

import (
	"context"
	"testing"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/stagedjob"
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
	defer terminate(ctx)

	// migrations up
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	// connect to database
	store, err := NewStore(dsn)
	require.NoError(t, err)

	t.Run("Create Staged Job", func(t *testing.T) {
		sj := &entity.StagedJob{
			JobName: stagedjob.JobNameSendReceipt,
			JobArgs: []byte("{\"data\": \"foo\"}"),
		}

		res, err := store.CreateStagedJob(ctx, sj)

		if assert.NoError(t, err) {
			sj.ID = res.ID
			assert.Equal(t, sj, res)
		}
	})
}
