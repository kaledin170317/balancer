package algoritms

import (
	"balancer/internal/domain/models"
)

type LeastConnections struct {
	backends []*models.Backend
}

func NewLeastConnections(backends []*models.Backend) *LeastConnections {
	return &LeastConnections{backends: backends}
}

func (lc *LeastConnections) Next() *models.Backend {
	var selected *models.Backend
	min := int64(-1)

	for _, b := range lc.backends {
		if !b.IsAlive() {
			continue
		}
		conns := b.ActiveConnections()
		if min == -1 || conns < min {
			min = conns
			selected = b
		}
	}

	return selected
}
