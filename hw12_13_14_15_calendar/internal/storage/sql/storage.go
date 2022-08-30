package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/config"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:golint
	"github.com/jmoiron/sqlx"
)

type dbStorageImpl struct {
	config config.DBConfig
	db     *sqlx.DB
}

func New(config config.DBConfig) *dbStorageImpl {
	return &dbStorageImpl{
		config: config,
	}
}

func (s *dbStorageImpl) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.config.Dsn)
	return err
}

func (s *dbStorageImpl) Close() error {
	return s.db.Close()
}

func (s *dbStorageImpl) CreateEvent(ctx context.Context, ev storage.Event) (storage.Event, error) {
	if err := checkIfAlreadyBooked(ctx, s, ev); err != nil {
		return storage.Event{}, err
	}

	res, err := s.db.NamedExecContext(ctx,
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

func checkIfAlreadyBooked(ctx context.Context, s *dbStorageImpl, ev storage.Event) error {
	var bookedEvents int
	err := s.db.GetContext(ctx, &bookedEvents,
		"select count(*) from events where start_at < $1 and end_at > $1", ev.StartTime)
	if err != nil {
		return err
	}
	if bookedEvents > 0 {
		return storage.ErrEventBooked
	}
	return nil
}

func (s *dbStorageImpl) Get(ctx context.Context, id int) (storage.Event, error) {
	event := storage.Event{}
	err := s.db.GetContext(ctx, &event, "select * from events where id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNotFound
		}
		return storage.Event{}, err
	}
	return event, nil
}

func (s *dbStorageImpl) Update(ctx context.Context, id int, event storage.Event) error {
	row, err := s.db.NamedQueryContext(ctx,
		"update events set title = :title, start_at = :start_at, end_at = :end_at where id = :id",
		map[string]interface{}{
			"id":       id,
			"title":    event.Title,
			"start_at": event.StartTime,
			"end_at":   event.StartTime.Add(event.Duration),
		})
	defer row.Close() //nolint
	return err
}

func (s *dbStorageImpl) Delete(ctx context.Context, id int) {
	s.db.QueryRowxContext(ctx, "delete from events where id = $1", id)
}

func (s *dbStorageImpl) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.SelectContext(ctx, &events, "select * from events where start_at > $1 and start_at < $2", start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}
