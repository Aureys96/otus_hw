package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
)

type StorageDao struct {
	s *Storage
}

func newEventDAO(s *Storage) storage.EventDao {
	return &StorageDao{s}
}

func (sd *StorageDao) CreateEvent(ctx context.Context, ev storage.Event) (storage.Event, error) {
	res, err := sd.s.db.NamedExecContext(ctx, "insert into events (title) values (:title) returning id",
		map[string]interface{}{
			"title": ev.Title,
		})
	if err != nil {
		return storage.Event{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ev, err
	}

	ev.ID = int(id)

	return ev, nil
}

func (sd *StorageDao) Get(ctx context.Context, id int) (storage.Event, error) {
	event := storage.Event{}
	err := sd.s.db.GetContext(ctx, &event, "select * from events where id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, errors.New("not found")
		}
		return storage.Event{}, err
	}
	return event, nil
}

func (sd *StorageDao) Update(ctx context.Context, id int, event storage.Event) error {
	_, err := sd.s.db.NamedQueryContext(ctx, "update events set title = :title where id = :id",
		map[string]interface{}{
			"id":    id,
			"title": event.Title,
		})
	return err
}

func (sd *StorageDao) Delete(ctx context.Context, id int) {
	sd.s.db.QueryRowxContext(ctx, "delete from events where id = $1", id)
}

func (sd *StorageDao) ListEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event
	err := sd.s.db.SelectContext(ctx, &events, "select * from events")
	if err != nil {
		return nil, err
	}
	return events, nil
}
