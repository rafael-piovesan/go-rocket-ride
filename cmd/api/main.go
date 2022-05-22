package main

import (
	"log"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/api/http"
	bunstore "github.com/rafael-piovesan/go-rocket-ride/datastore/bun"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/stripemock"
)

func main() {
	// app's config values
	cfg, err := rocketride.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// database connection
	store, err := bunstore.NewStore(cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	// Replace the original Stripe API Backend with its mock
	stripemock.Init()

	// http server
	httpServer := http.NewServer(cfg, store)
	httpServer.Start()
}
