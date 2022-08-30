package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
)

type inmemoryStorageImpl struct {
	mu   sync.RWMutex
	data map[int]storage.Event
}

func (s *inmemoryStorageImpl) Connect(_ context.Context) error {
	return nil
}

func (s *inmemoryStorageImpl) Close() error {
	return nil
}

func New() *inmemoryStorageImpl {
	return &inmemoryStorageImpl{data: make(map[int]storage.Event)}
}

func (s *inmemoryStorageImpl) CreateEvent(_ context.Context, event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ev := range s.data {
		if ev.StartTime.UnixNano() < event.StartTime.UnixNano() &&
			ev.StartTime.Add(ev.Duration).UnixNano() > event.StartTime.UnixNano() {
			return storage.Event{}, storage.ErrEventBooked
		}
	}

	s.data[event.ID] = event
	return event, nil
}

func (s *inmemoryStorageImpl) Get(_ context.Context, id int) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return storage.Event{}, storage.ErrNotFound
	}
	return s.data[id], nil
}

func (s *inmemoryStorageImpl) Update(_ context.Context, id int, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return storage.ErrNotFound
	}

	s.data[id] = event
	return nil
}

func (s *inmemoryStorageImpl) Delete(_ context.Context, id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
}

func (s *inmemoryStorageImpl) ListEvents(_ context.Context, start, end time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]storage.Event, len(s.data))
	for i, event := range s.data {
		if event.StartTime.After(start) && event.StartTime.Add(event.Duration).Before(end) {
			events[i] = event
		}
	}
	return events, nil
}
