package main

import (
	"database/sql"
	"log"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/adapters/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/api/http"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	// app's config values
	cfg, err := rocketride.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// database connection
	dsn := cfg.DBSource
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	store := datastore.NewStore(db)

	// http server
	httpServer := http.NewServer(cfg, store)
	httpServer.Start()
}
