package cheker

import (
	"balancer/internal/domain/models"
	"balancer/internal/logger"
	"context"
	"net"
	"net/url"
	"time"
)

type Checker struct {
	interval time.Duration
	backends []*models.Backend
}

func NewChecker(backends []*models.Backend, interval time.Duration) *Checker {
	return &Checker{
		backends: backends,
		interval: interval,
	}
}

func (hc *Checker) StartCheck(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("checker started", "interval", hc.interval)

	go func() {
		ticker := time.NewTicker(hc.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("checker stopped")
				return
			case <-ticker.C:
				for _, b := range hc.backends {
					if !b.IsAlive() {
						go hc.checkTCP(ctx, b)
					}
				}
			}
		}
	}()
}

func (hc *Checker) checkTCP(ctx context.Context, b *models.Backend) {
	log := logger.FromContext(ctx)

	if isBackendAlive(b.URL) {
		if b.SetAlive(true) {
			log.Info("backend marked alive", "url", b.URL.String())
		}
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		logger.Base().Warn("backend unreachable", "host", u.Host, "err", err)
		return false
	}
	_ = conn.Close()
	return true
}
