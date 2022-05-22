package main

import (
	"log"

	rocketride "github.com/rafael-piovesan/go-rocket-ride/v2"
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/http"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/stripemock"
)

func main() {
	// app's config values
	cfg, err := rocketride.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// database connection
	store, err := datastore.NewStore(cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	// Replace the original Stripe API Backend with its mock
	stripemock.Init()

	// http server
	httpServer := http.NewServer(cfg, store)
	httpServer.Start()
}
