package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadFromEnv() (*Config, error) {
	var cfg Config
	cfg.ListenAddr = get("LISTEN_ADDR", ":8080")
	cfg.Algorithm = get("ALGORITHM", "round-robin")
	cfg.Backends = strings.Split(get("BACKENDS", "https://httpbin.org/get"), ",")
	if len(cfg.Backends) == 1 && cfg.Backends[0] == "" {
		cfg.Backends = nil
	}

	var err error
	if cfg.RateLimit.Capacity, err = getInt("RATE_CAPACITY", 100); err != nil {
		return nil, err
	}
	if cfg.RateLimit.RefillRate, err = getInt("RATE_REFILL", 10); err != nil {
		return nil, err
	}
	if cfg.HealthCheck.Interval, err = getDuration("HC_INTERVAL", 5*time.Second); err != nil {
		return nil, err
	}
	if cfg.HealthCheck.Timeout, err = getDuration("HC_TIMEOUT", 2*time.Second); err != nil {
		return nil, err
	}

	cfg.DatabaseDSN = get("DB_DSN", "postgres://postgres:password@localhost:5555/balancer?sslmode=disable")
	cfg.RateLimit.Window = time.Second / time.Duration(max(cfg.RateLimit.RefillRate, 1))

	return &cfg, nil
}

func get(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getInt(key string, def int) (int, error) {
	v := get(key, "")
	if v == "" {
		return def, nil
	}
	return strconv.Atoi(v)
}

func getDuration(key string, def time.Duration) (time.Duration, error) {
	v := get(key, "")
	if v == "" {
		return def, nil
	}
	return time.ParseDuration(v)
}
