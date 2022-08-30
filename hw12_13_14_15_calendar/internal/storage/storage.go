package storage

import (
	"context"
	"time"
)

type IStorage interface {
	CreateEvent(ctx context.Context, event Event) (Event, error)
	Get(ctx context.Context, id int) (Event, error)
	Update(ctx context.Context, id int, event Event) error
	Delete(ctx context.Context, id int)
	ListEvents(ctx context.Context, start, end time.Time) ([]Event, error)
	Connect(ctx context.Context) error
	Close() error
}
