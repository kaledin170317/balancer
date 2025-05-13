package ratelimit

import (
	"balancer/internal/domain/models"
	"balancer/internal/domain/usecases"
	"balancer/internal/logger"
	"context"
	"sync"
)

type Limiter struct {
	buckets sync.Map // map[string]*TokenBucket
	uc      *usecases.ClientUseCase
}

func NewLimiterFromUseCase(uc *usecases.ClientUseCase) *Limiter {
	return &Limiter{uc: uc}
}

func (l *Limiter) getOrCreateBucket(clientID string) *TokenBucket {
	if bucket, ok := l.buckets.Load(clientID); ok {
		return bucket.(*TokenBucket)
	}

	client, err := l.uc.Get(context.Background(), clientID)
	if err != nil {
		logger.Warn(nil, "unknown client, using default limiter", "clientID", clientID)
		client = &models.Client{
			ID:         clientID,
			Capacity:   100,
			RatePerSec: 10,
		}
	}

	b := NewTokenBucket(client.Capacity, client.RatePerSec)
	l.buckets.Store(clientID, b)
	return b
}

func (l *Limiter) Allow(clientID string) bool {
	return l.getOrCreateBucket(clientID).Allow()
}
