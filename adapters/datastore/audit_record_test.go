//go:build integration
// +build integration

package datastore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/audit"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/migrate"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testcontainer"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tabbed/pqtype"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestAuditRecord(t *testing.T) {
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
	err = testfixtures.Load(dsn, []string{"db/fixtures/users"}, map[string]interface{}{
		"UserId":    userID,
		"UserEmail": gofakeit.Email(),
	})
	require.NoError(t, err)

	// conntect to database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	store := NewStore(db)

	t.Run("Create Audit Record", func(t *testing.T) {
		ip := pqtype.CIDR{}
		err := ip.Scan(gofakeit.IPv4Address())
		require.NoError(t, err)

		ar := &entity.AuditRecord{
			Action:       audit.ActionCreateRide,
			Data:         []byte("{\"data\": \"foo\"}"),
			OriginIP:     ip.IPNet.String(),
			ResourceID:   int64(gofakeit.Number(0, 1000)),
			ResourceType: audit.ResourceTypeRide,
			UserID:       userID,
		}

		res, err := store.CreateAuditRecord(ctx, ar)

		if assert.NoError(t, err) {
			ar.ID = res.ID
			assert.Equal(t, ar, res)
		}
	})
}
