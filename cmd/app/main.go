package main

import (
	"log"

	"github.com/captain-corgi/go-ohlc-history-service/config"
	"github.com/captain-corgi/go-ohlc-history-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
