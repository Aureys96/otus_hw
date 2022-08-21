package memorystorage

import (
	"context"
	"errors"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
)

func newEventDAO(st *Storage) *StorageDao {
	return &StorageDao{st}
}

type StorageDao struct {
	storage *Storage
}

func (s *StorageDao) CreateEvent(_ context.Context, event storage.Event) (storage.Event, error) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.data[event.ID] = event
	return event, nil
}

func (s *StorageDao) Get(_ context.Context, id int) (storage.Event, error) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	if _, ok := s.storage.data[id]; !ok {
		return storage.Event{}, errors.New("event not found")
	}
	return s.storage.data[id], nil
}

func (s *StorageDao) Update(_ context.Context, id int, event storage.Event) error {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	_, ok := s.storage.data[id]
	if !ok {
		return errors.New("event for update not found")
	}
	s.storage.data[id] = event
	return nil
}

func (s *StorageDao) Delete(_ context.Context, id int) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	delete(s.storage.data, id)
}

func (s *StorageDao) ListEvents(_ context.Context) ([]storage.Event, error) {
	//TODO implement me
	panic("implement me")
}
