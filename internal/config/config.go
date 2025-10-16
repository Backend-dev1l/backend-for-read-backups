package config

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Logger     Logger     `envPrefix:"LOGGER_"`
	HTTP       HTTP       `envPrefix:"HTTP_"`
	Postgres   Postgres   `envPrefix:"POSTGRES_"`
	App        AppConfig  `envPrefix:"APP_"`
	PoolConfig PoolConfig `envPrefix:"POOL_"`
	TimeOuts   TimeOuts   `envPrefix:"TIMEOUTS_"`
}

type HTTP struct {
	Port int    `env:"PORT" envDefault:"8080" validate:"min=1024,max=65535"`
	Host string `env:"HOST" envDefault:"localhost" validate:"hostname|ip"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
}

type Postgres struct {
	Host     string `env:"HOST"     envDefault:"localhost" validate:"required,hostname|ip"`
	Port     int    `env:"PORT"     envDefault:"5432"      validate:"required,min=1024,max=65535"`
	User     string `env:"USER"     envDefault:"postgres"  validate:"required"`
	Password string `env:"PASSWORD" envDefault:""          validate:"required"`
	DBName   string `env:"DBNAME"   envDefault:"postgres"      validate:"required"`
	SSLMode  string `env:"SSLMODE"  envDefault:"disable"   validate:"oneof=disable require verify-ca verify-full"`
}

type TimeOuts struct {
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s" validate:"min=1s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s" validate:"min=1s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"10s" validate:"min=1s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s" validate:"min=1s"`
	PerRequestTimeout time.Duration `env:"PER_REQUEST_TIMEOUT" envDefault:"5s" validate:"min=1s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s" validate:"min=1s"`
}

type AppConfig struct {
	Service string `env:"SERVICE" envDefault:"test-http" validate:"required"`
	Env     string `env:"ENV"     envDefault:"development" validate:"required,oneof=development staging production"`
	Version string `env:"VERSION" envDefault:"1.0.0" validate:"required"`
}

type PoolConfig struct {
	MaxConns          int32         `env:"MAX_CONNS" envDefault:"16" validate:"min=1,max=100"`
	MinConns          int32         `env:"MIN_CONNS" envDefault:"4" validate:"min=1,max=100"`
	MaxConnLifetime   time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"1h" validate:"min=1m"`
	MaxConnIdleTime   time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"15m" validate:"min=1m"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"1m" validate:"min=10s"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.HTTP.Host, strconv.Itoa(c.HTTP.Port))
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
