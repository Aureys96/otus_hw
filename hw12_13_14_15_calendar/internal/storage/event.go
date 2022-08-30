package storage

import "time"

type Event struct {
	ID        int
	Title     string
	StartTime time.Time `db:"start_at"`
	Duration  time.Duration
	UserID    int `db:"user_id"`
}
