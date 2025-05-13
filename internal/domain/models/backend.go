package models

import (
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL         *url.URL
	alive       atomic.Bool
	activeConns atomic.Int64
}

func NewBackend(u *url.URL) *Backend {
	b := &Backend{URL: u}
	b.alive.Store(true)
	return b
}

func (b *Backend) SetAlive(alive bool) bool {
	previous := b.alive.Load()
	b.alive.Store(alive)
	return previous != alive
}

func (b *Backend) IsAlive() bool {
	return b.alive.Load()
}

func (b *Backend) ActiveConnections() int64 {
	return b.activeConns.Load()
}

func (b *Backend) IncConnections() {
	b.activeConns.Add(1)
}

func (b *Backend) DecConnections() {
	b.activeConns.Add(-1)
}
