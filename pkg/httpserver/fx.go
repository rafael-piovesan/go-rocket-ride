package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"go.uber.org/fx"
)

func Invoke(lc fx.Lifecycle, e *echo.Echo, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
			// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
			// resources to be properly freed. If it was the other way around, the application would just
			// exit without gracefully shutting down the server.
			// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
			go func() {
				if err := e.Start(cfg.ServerAddress); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				log.Errorf("error shutting down server: %v", err)
			} else {
				log.Info("server shutdown gracefully")
			}
			return nil
		},
	})
}
