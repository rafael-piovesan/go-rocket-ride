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
	"github.com/stretchr/testify/require"
)

func TestStagedJobs(t *testing.T) {
	ctx := context.Background()

	// create a new Postgres container to be used in tests
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer terminate(ctx)

	// run all migrations
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	testQueries := New(db)

	t.Run("Create Staged Job", func(t *testing.T) {
		j := CreateStagedJobParams{
			JobName: gofakeit.RandomString([]string{"completer", "enqueuer", "reaper"}),
			JobArgs: []byte("{}"),
		}
		r, err := testQueries.CreateStagedJob(ctx, j)
		require.NoError(t, err)

		require.Equal(t, j.JobName, r.JobName)
		require.Equal(t, j.JobArgs, r.JobArgs)
	})
}
