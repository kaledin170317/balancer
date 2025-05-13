package algoritms

import (
	"balancer/internal/domain/models"
	"math/rand"
	"sync"
	"time"
)

type Random struct {
	backends []*models.Backend
	mu       sync.Mutex
	rng      *rand.Rand
}

func NewRandom(backends []*models.Backend) *Random {
	return &Random{
		backends: backends,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) Next() *models.Backend {
	r.mu.Lock()
	defer r.mu.Unlock()

	alive := make([]*models.Backend, 0)
	for _, b := range r.backends {
		if b.IsAlive() {
			alive = append(alive, b)
		}
	}

	if len(alive) == 0 {
		return nil
	}

	return alive[r.rng.Intn(len(alive))]
}
