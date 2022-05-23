package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/http/handler"
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/http/middleware"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/usecase"
	"go.uber.org/fx"
)

func routes(e *echo.Echo, userStore datastore.User, ride handler.Ride) {
	e.Use(middleware.NewOriginIP())
	e.Use(middleware.NewErrorMapper())
	e.Use(middleware.NewUser(userStore))

	// Routes
	e.POST("/", ride.Create)
}

var Module = fx.Options(
	fx.Provide(
		usecase.NewRide,
		handler.NewRide,
	),
	fx.Invoke(routes),
)
