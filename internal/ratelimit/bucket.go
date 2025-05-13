package ratelimit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	capacity   int
	ratePerSec int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		ratePerSec: rate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	newTokens := int(elapsed * float64(tb.ratePerSec))
	if newTokens > 0 {
		tb.tokens = min(tb.tokens+newTokens, tb.capacity)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
