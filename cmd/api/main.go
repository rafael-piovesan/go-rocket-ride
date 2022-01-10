package main

import (
	"database/sql"
	"log"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/adapters/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/api/http"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/stripemock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	// app's config values
	cfg, err := rocketride.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// database connection
	dsn := cfg.DBSource
	sqldb, err := sql.Open("pg", dsn)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	store := datastore.NewStore(db)

	// Replace the original Stripe API Backend with its mock
	stripemock.Init()

	// http server
	httpServer := http.NewServer(cfg, store)
	httpServer.Start()
}
