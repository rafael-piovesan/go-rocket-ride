//go:build integration
// +build integration

package bun

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
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
	userEmail := gofakeit.Email()
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": userEmail,
	})
	require.NoError(t, err)

	// conntect to database
	store, err := NewStore(dsn)
	require.NoError(t, err)

	t.Run("User not found", func(t *testing.T) {
		_, err := store.GetUserByEmail(ctx, gofakeit.FarmAnimal())
		assert.ErrorIs(t, err, entity.ErrNotFound)
	})

	t.Run("User found", func(t *testing.T) {
		u, err := store.GetUserByEmail(ctx, userEmail)
		assert.NoError(t, err)
		assert.Equal(t, userID, u.ID)
		assert.Equal(t, userEmail, u.Email)
	})
}
