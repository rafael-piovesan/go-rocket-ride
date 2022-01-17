package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
	im := cstmiddleware.NewIPMiddleware()
	e.Use(im.Handle)

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
	// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
	// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
	// resources to be properly freed. If it was the other way around, the application would just
	// exit without gracefully shutting down the server.
	// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
	go func() {
		if err := s.e.Start(s.cfg.ServerAddress); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error running server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down server: %v", err)
	} else {
		log.Info("server shutdown gracefully")
	}
}
