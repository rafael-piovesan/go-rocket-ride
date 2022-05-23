//go:build integration
// +build integration

package repo

import (
	"context"
	"testing"

	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/testcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

type book struct {
	ID     int64
	Title  string
	Author string
}

func TestCRUDRepository(t *testing.T) {
	ctx := context.Background()

	dsn, terminate, err := testcontainer.NewPostgresContainer()
	require.NoError(t, err)
	defer func() { _ = terminate(ctx) }()

	db, err := db.Connect(config.Config{DBSource: dsn})
	require.NoError(t, err)

	_, err = db.NewCreateTable().Model(&book{}).Exec(ctx)
	require.NoError(t, err)

	repo := New[book](db)
	books := []book{
		{Title: "foo1", Author: "bar1"},
		{Title: "foo2", Author: "bar2"},
	}

	t.Run("save model", func(t *testing.T) {
		bks, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(bks))

		for i := range books {
			err = repo.Save(ctx, &books[i])
			assert.NoError(t, err)
		}
	})

	t.Run("find all", func(t *testing.T) {
		bks, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, books, bks)
	})

	t.Run("find all with criteria", func(t *testing.T) {
		c := func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("title = ?", books[0].Title)
		}

		bks, err := repo.FindAll(ctx, c)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bks))
		assert.Equal(t, books[0], bks[0])
	})

	t.Run("find one", func(t *testing.T) {
		b, err := repo.FindOne(ctx)
		assert.NoError(t, err)
		assert.Equal(t, books[0], b)
	})

	t.Run("find one with criteria", func(t *testing.T) {
		c := func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("title = ?", books[1].Title)
		}

		b, err := repo.FindOne(ctx, c)
		assert.NoError(t, err)
		assert.Equal(t, books[1], b)
	})

	t.Run("delete model", func(t *testing.T) {
		err = repo.Delete(ctx, &books[0])
		assert.NoError(t, err)

		bks, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bks))
		assert.Equal(t, books[1], bks[0])
	})

	t.Run("update model", func(t *testing.T) {
		books[1].Title = "foo3"
		err = repo.Update(ctx, &books[1])
		assert.NoError(t, err)

		bks, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bks))
		assert.Equal(t, books[1], bks[0])
	})
}
