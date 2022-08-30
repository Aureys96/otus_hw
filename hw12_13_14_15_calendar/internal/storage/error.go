package storage

import "errors"

var (
	ErrNotFound    = errors.New("event not found")
	ErrEventBooked = errors.New("event already booked")
)
