// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/captain-corgi/go-ohlc-history-service/config"
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	// Swagger docs.
	_ "github.com/captain-corgi/go-ohlc-history-service/docs"
	"github.com/captain-corgi/go-ohlc-history-service/internal/usecase"
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
)

// NewEchoRouter -.
// Swagger spec:
// @title       OHLC Handler
// @description Use for import/export OHLC data.
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewEchoRouter(handler *echo.Echo, l logger.Interface, t usecase.OHLCUseCase, cfg *config.Config) {
	// Options
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	// Swagger
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	// K8s probe
	handler.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{"OK"})
	})

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(handler)

	// Routers
	h := handler.Group("/v1")
	{
		newOHLCRoutes(h, t, l, cfg)
	}
}
