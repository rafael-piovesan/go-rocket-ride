package main

import (
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/http"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/db"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/httpserver"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/stripemock"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			db.Connect,
			db.ConnectionHandle,
			datastore.New,
			datastore.NewIdempotencyKey,
			datastore.NewUser,
			httpserver.New,
		),
		// Loading HTTP routes & handlers
		http.Module,
		// Replace the original Stripe API Backend with its mock
		fx.Invoke(stripemock.Init),
		fx.Invoke(httpserver.Invoke),
	).Run()
}
