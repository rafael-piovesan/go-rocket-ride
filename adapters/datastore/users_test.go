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

func TestUsers(t *testing.T) {
	ctx := context.Background()

	// create a new Postgres container to be used in tests
	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer terminate(ctx)

	// run all migrations
	err = migrate.Up(dsn, "db/migrations")
	require.NoError(t, err)

	userID := int64(gofakeit.Number(0, 10000))
	email := gofakeit.Email()

	// load database test fixtures separately
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": email,
	})
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	testQueries := New(db)

	// all set, finally run the tests
	t.Run("Get User by ID", func(t *testing.T) {
		u, err := testQueries.GetUserByID(ctx, userID)
		require.NoError(t, err)

		require.Equal(t, userID, u.ID)
		require.Equal(t, email, u.Email)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		u, err := testQueries.GetUserByEmail(ctx, email)
		require.NoError(t, err)

		require.Equal(t, userID, u.ID)
		require.Equal(t, email, u.Email)
	})
}
