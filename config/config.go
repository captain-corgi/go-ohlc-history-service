package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		MySQL `yaml:"mysql"`
		App   `yaml:"app"`
		Log   `yaml:"logger"`
		HTTP  `yaml:"http"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		// Port: '8080'
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		// Host: '0.0.0.0'
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		// ReadTimeout: '60'
		ReadTimeout int `env-required:"true" yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
		// WriteTimeout: '60'
		WriteTimeout int `env-required:"true" yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
		// ShutdownTimeout: '60'
		ShutdownTimeout int `env-required:"true" yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
		// ProcessTimeout process_timeout: '60'
		ProcessTimeout int `env-required:"true" yaml:"process_timeout" env:"HTTP_PROCESS_TIMEOUT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// MySQL -.
	MySQL struct {
		URL             string `env-required:"true"                           env:"MYSQL_URL"`
		USER            string `env-required:"true"                           env:"MYSQL_USER"`
		PASS            string `env-required:"true"                           env:"MYSQL_PASS"`
		SCHEMA          string `env-required:"true"                           env:"MYSQL_SCHEMA"`
		MaxOpenConns    int    `env-required:"true" yaml:"max_open_conns"     env:"MYSQL_MAX_OPEN_CONNS"`
		MaxIdleConns    int    `env-required:"true" yaml:"max_idle_conns"     env:"MYSQL_MAX_IDLE_CONNS"`
		ConnMaxLifetime int    `env-required:"true" yaml:"conn_max_lifetime"  env:"MYSQL_CONN_MAX_LIFETIME"`
		ConnPingTimeout int    `env-required:"true" yaml:"conn_ping_timeout"  env:"MYSQL_CONN_PING_TIMEOUT"`
		WriteThreshold  int    `env-required:"true" yaml:"write_threshold"    env:"MYSQL_WRITE_THRESHOLD"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
