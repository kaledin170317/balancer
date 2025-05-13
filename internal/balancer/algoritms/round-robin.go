package algoritms

import (
	"balancer/internal/domain/models"
	"sync/atomic"
)

type RoundRobin struct {
	backends []*models.Backend
	counter  uint64
}

func NewRoundRobin(backends []*models.Backend) *RoundRobin {
	return &RoundRobin{backends: backends}
}

func (r *RoundRobin) Next() *models.Backend {
	n := len(r.backends)
	if n == 0 {
		return nil
	}

	start := atomic.AddUint64(&r.counter, 1)
	for i := 0; i < n; i++ {
		idx := int((start + uint64(i)) % uint64(n))
		b := r.backends[idx]
		if b.IsAlive() {
			return b
		}
	}
	
	return nil
}
