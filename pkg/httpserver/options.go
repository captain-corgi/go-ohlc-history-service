package httpserver

import (
	"github.com/captain-corgi/go-ohlc-history-service/pkg/logger"
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Logger -.
func Logger(logger logger.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout int) Option {
	return func(s *Server) {
		s.server.ReadTimeout = time.Duration(timeout) * time.Second
	}
}

// WriteTimeout -.
func WriteTimeout(timeout int) Option {
	return func(s *Server) {
		s.server.WriteTimeout = time.Duration(timeout) * time.Second
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout int) Option {
	return func(s *Server) {
		s.shutdownTimeout = time.Duration(timeout) * time.Second
	}
}
