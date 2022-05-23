package db

import (
	"database/sql"

	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	// loading bun's official Postgres driver.
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func Connect(cfg config.Config) (*bun.DB, error) {
	sqldb, err := sql.Open("pg", cfg.DBSource)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}

func ConnectionHandle(db *bun.DB) bun.IDB {
	return db
}
