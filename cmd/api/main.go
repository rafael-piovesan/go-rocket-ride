package main

import (
	"log"
	"strings"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/api/http"
	bunstore "github.com/rafael-piovesan/go-rocket-ride/datastore/bun"
	sqlcstore "github.com/rafael-piovesan/go-rocket-ride/datastore/sqlc"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/stripemock"
)

func main() {
	// app's config values
	cfg, err := rocketride.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// database connection
	da := strings.ToLower(cfg.DatastoreAccess)
	var store rocketride.Datastore

	switch da {
	case "bun":
		store, err = bunstore.NewStore(cfg.DBSource)
	case "sqlc":
		store, err = sqlcstore.NewStore(cfg.DBSource)
	default:
		log.Printf("invalid store type %v", cfg.DatastoreAccess)
		store, err = bunstore.NewStore(cfg.DBSource)
	}

	log.Printf("using [%v] datastore access", da)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	// Replace the original Stripe API Backend with its mock
	stripemock.Init()

	// http server
	httpServer := http.NewServer(cfg, store)
	httpServer.Start()
}
