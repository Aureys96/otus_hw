package memorystorage

import (
	"context"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
)

type StorageDao struct {
	storage *Storage
}

func newEventDAO(st *Storage) *StorageDao {
	return &StorageDao{st}
}

func (s *StorageDao) CreateEvent(_ context.Context, event storage.Event) (storage.Event, error) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	for _, ev := range s.storage.data {
		if ev.StartTime.UnixNano() < event.StartTime.UnixNano() &&
			ev.StartTime.Add(ev.Duration).UnixNano() > event.StartTime.UnixNano() {
			return storage.Event{}, storage.ErrEventBooked
		}
	}

	s.storage.data[event.ID] = event
	return event, nil
}

func (s *StorageDao) Get(_ context.Context, id int) (storage.Event, error) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	if _, ok := s.storage.data[id]; !ok {
		return storage.Event{}, storage.ErrNotFound
	}
	return s.storage.data[id], nil
}

func (s *StorageDao) Update(_ context.Context, id int, event storage.Event) error {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	if _, ok := s.storage.data[id]; !ok {
		return storage.ErrNotFound
	}

	s.storage.data[id] = event
	return nil
}

func (s *StorageDao) Delete(_ context.Context, id int) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	delete(s.storage.data, id)
}

func (s *StorageDao) ListEvents(_ context.Context, start, end time.Time) ([]storage.Event, error) {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	events := make([]storage.Event, len(s.storage.data))
	for i, event := range s.storage.data {
		if event.StartTime.After(start) && event.StartTime.Add(event.Duration).Before(end) {
			events[i] = event
		}
	}
	return events, nil
}
