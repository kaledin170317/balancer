package config

import "time"

type Config struct {
	ListenAddr  string            `json:"listen_addr" yaml:"listen_addr" env:"LISTEN_ADDR"`
	Backends    []string          `json:"backends" yaml:"backends" env:"BACKENDS"`
	Algorithm   string            `json:"algorithm" yaml:"algorithm" env:"ALGORITHM"`
	RateLimit   RateLimitConfig   `json:"rate_limit" yaml:"rate_limit"`
	HealthCheck HealthCheckConfig `json:"health_check" yaml:"health_check"`
	DatabaseDSN string            `json:"database_dsn" yaml:"database_dsn" env:"DB_DSN"`
}

type RateLimitConfig struct {
	Capacity   int           `json:"capacity"     yaml:"capacity"     env:"RATE_CAPACITY"`
	RefillRate int           `json:"refill_rate"  yaml:"refill_rate"  env:"RATE_REFILL"`
	Window     time.Duration `json:"-"            yaml:"-"`
}

type HealthCheckConfig struct {
	Interval time.Duration `json:"interval" yaml:"interval" env:"HC_INTERVAL"`
	Timeout  time.Duration `json:"timeout"  yaml:"timeout"  env:"HC_TIMEOUT"`
}
