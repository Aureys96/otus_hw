package memorystorage

import (
	"context"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[int]storage.Event
	dao  storage.EventDao
}

func (s *Storage) Connect(_ context.Context) error {
	panic("unsupported")
}

func (s *Storage) Close() error {
	panic("unsupported")
}

func New() *Storage {
	return &Storage{data: make(map[int]storage.Event)}
}

func (s *Storage) DAO() storage.EventDao {
	if s.dao == nil {
		s.dao = newEventDAO(s)
	}
	return s.dao
}
