// Package mysql implements mysql connection.
package mysql

import (
	"context"
	"database/sql"
	"github.com/captain-corgi/go-ohlc-history-service/config"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// MySQL -.
type MySQL struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *sql.DB
}

// New -.
func New(config config.MySQL, opts ...Option) (*MySQL, error) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   config.USER,
		Passwd: config.PASS,
		Net:    "tcp",
		Addr:   config.URL,
		DBName: config.SCHEMA,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Minute)

	// Create a new MySQL instance.
	msql := &MySQL{
		DB: db,
	}

	// Custom options
	for _, opt := range opts {
		opt(msql)
	}

	// Check connection.
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(config.ConnPingTimeout)*time.Second)
	defer cancelFunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}

	return msql, nil
}

// Close -.
func (p *MySQL) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}
