//go:build integration
// +build integration

package datastore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/require"
	"github.com/tabbed/pqtype"
)

func TestAuditRecord(t *testing.T) {
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
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": gofakeit.Email(),
	})
	require.NoError(t, err)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	testQueries := New(db)

	t.Run("Create Audit Record", func(t *testing.T) {
		ip := pqtype.CIDR{}
		err := ip.Scan(gofakeit.IPv4Address())
		require.NoError(t, err)
		require.True(t, ip.Valid)

		a := CreateAuditRecordParams{
			Action:       gofakeit.RandomString([]string{"create", "update", "delete"}),
			CreatedAt:    time.Now(),
			Data:         []byte("{}"),
			OriginIp:     ip.IPNet.String(),
			ResourceID:   int64(gofakeit.Number(0, 10000)),
			ResourceType: gofakeit.RandomString([]string{"user", "idempotency_key", "ride"}),
			UserID:       userID,
		}

		r, err := testQueries.CreateAuditRecord(ctx, a)
		require.NoError(t, err)

		require.Equal(t, a.Action, r.Action)
		require.Equal(t, a.Data, r.Data)
		require.Equal(t, a.OriginIp, r.OriginIp)
		require.Equal(t, a.ResourceID, r.ResourceID)
		require.Equal(t, a.ResourceType, r.ResourceType)
		require.Equal(t, a.UserID, r.UserID)
	})
}
