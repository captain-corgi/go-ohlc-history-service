// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/cache"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"syscall"

	"github.com/captain-corgi/go-ohlc-history-service/config"
	v1 "github.com/captain-corgi/go-ohlc-history-service/internal/controller/http/v1"
	"github.com/captain-corgi/go-ohlc-history-service/internal/usecase"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/httpserver"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/mysql"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Info("Starting application...")

	cache.InMemory = cache.NewInMemoryDB()

	// Repository
	mySQL, err := mysql.New(cfg.MySQL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - mysql.New: %w", err))
	}
	defer mySQL.Close()

	// Use case
	ohlcUseCase := usecase.NewOHLCUseCase(mySQL.DB, *l)

	// HTTP Server
	echoServer := echo.New()
	v1.NewEchoRouter(echoServer, l, *ohlcUseCase, cfg)
	httpServer := httpserver.New(echoServer,
		httpserver.Logger(*l),
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.ReadTimeout),
		httpserver.WriteTimeout(cfg.HTTP.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
