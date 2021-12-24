package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/api/http/handler"
	cstmiddleware "github.com/rafael-piovesan/go-rocket-ride/api/http/middleware"
	"github.com/rafael-piovesan/go-rocket-ride/usecase"
)

type Server struct {
	e   *echo.Echo
	cfg rocketride.Config
}

func NewServer(cfg rocketride.Config, store rocketride.Datastore) *Server {
	e := echo.New()

	// Middleware
	um := cstmiddleware.NewUserMiddleware(store)
	e.Use(um.Handle)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers
	ride := usecase.NewRideUseCase(cfg, store)
	rideHandler := handler.NewRideHandler(ride)

	// Routes
	e.POST("/", rideHandler.Create)

	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) Start() {
	// Start server
	go func() {
		if err := s.e.Start(s.cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Fatal(err)
	}
}
