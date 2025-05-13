package config

import (
	"flag"
	"strings"
	"time"
)

func ParseFromFlags() *Config {
	var (
		listenAddr = flag.String("listen", "", "Address to listen on (e.g. :8080)")
		backends   = flag.String("backends", "", "Comma-separated list of backend URLs")
		algorithm  = flag.String("algorithm", "", "Load balancing algorithm: round-robin, least-connections, random")
		dbDSN      = flag.String("db-dsn", "", "DSN for persistent client storage")

		rateCapacity = flag.Int("rate-capacity", 0, "Token bucket capacity")
		rateRefill   = flag.Int("rate-refill", 0, "Token refill rate (tokens/sec)")

		hcInterval = flag.Duration("hc-interval", 0, "Health check interval (e.g. 5s, 1m)")
		hcTimeout  = flag.Duration("hc-timeout", 0, "Health check timeout")
	)

	flag.Parse()

	cfg := &Config{
		ListenAddr:  *listenAddr,
		Algorithm:   *algorithm,
		DatabaseDSN: *dbDSN,
	}

	if *backends != "" {
		parts := strings.Split(*backends, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				cfg.Backends = append(cfg.Backends, trimmed)
			}
		}
	}

	cfg.RateLimit.Capacity = *rateCapacity
	cfg.RateLimit.RefillRate = *rateRefill
	if *rateRefill > 0 {
		cfg.RateLimit.Window = time.Second / time.Duration(*rateRefill)
	}

	cfg.HealthCheck.Interval = *hcInterval
	cfg.HealthCheck.Timeout = *hcTimeout

	return cfg
}
