//go:build integration
// +build integration

package datastore

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
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
	userEmail := gofakeit.Email()
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": userEmail,
	})
	require.NoError(t, err)

	// conntect to database
	db, _ := db.Connect(config.Config{DBSource: dsn})
	store := NewUser(db)

	t.Run("User not found", func(t *testing.T) {
		_, err := store.FindOne(ctx, UserWithEmail(gofakeit.FarmAnimal()))
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("User found", func(t *testing.T) {
		u, err := store.FindOne(ctx, UserWithEmail(userEmail))
		assert.NoError(t, err)
		assert.Equal(t, userID, u.ID)
		assert.Equal(t, userEmail, u.Email)
	})
}
