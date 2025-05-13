package models

import "time"

type Client struct {
	ID         string    `db:"id"`
	Capacity   int       `db:"capacity"`
	RatePerSec int       `db:"rate_per_sec"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
