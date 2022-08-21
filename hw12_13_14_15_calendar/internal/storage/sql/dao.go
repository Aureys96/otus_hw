package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
)

type StorageDao struct {
	s *Storage
}

func newEventDAO(s *Storage) storage.EventDao {
	return &StorageDao{s}
}

func (sd *StorageDao) CreateEvent(ctx context.Context, ev storage.Event) (storage.Event, error) {
	res, err := sd.s.db.NamedExecContext(ctx,
		`insert into events (title, start_at, end_at, user_id) 
                  values (:title, :start_at, :end_at, :uid) returning id`,
		map[string]interface{}{
			"title":    ev.Title,
			"start_at": ev.StartTime,
			"end_at":   ev.StartTime.Add(ev.Duration),
			"uid":      ev.UserID,
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
			return storage.Event{}, storage.ErrNotFound
		}
		return storage.Event{}, err
	}
	return event, nil
}

func (sd *StorageDao) Update(ctx context.Context, id int, event storage.Event) error {
	_, err := sd.s.db.NamedQueryContext(ctx,
		"update events set title = :title, start_at = :start_at, end_at = :end_at where id = :id",
		map[string]interface{}{
			"id":       id,
			"title":    event.Title,
			"start_at": event.StartTime,
			"end_at":   event.StartTime.Add(event.Duration),
		})
	return err
}

func (sd *StorageDao) Delete(ctx context.Context, id int) {
	sd.s.db.QueryRowxContext(ctx, "delete from events where id = $1", id)
}

func (sd *StorageDao) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := sd.s.db.SelectContext(ctx, &events, "select * from events where start_at > $1 and start_at < $2", start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}
