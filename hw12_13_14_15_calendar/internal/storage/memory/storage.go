package memorystorage

import (
	"context"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func (s Storage) Events() storage.IEventStorage {
	//TODO implement me
	panic("implement me")
}

func (s Storage) Connect(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) Close() error {
	//TODO implement me
	panic("implement me")
}

func New() *Storage {
	return &Storage{}
}

// TODO
