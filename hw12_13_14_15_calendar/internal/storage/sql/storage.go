package sqlstorage

import (
	"context"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/config"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:golint
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	config config.DBConfig
	db     *sqlx.DB
	dao    storage.EventDao
}

func (s *Storage) DAO() storage.EventDao {
	if s.dao == nil {
		s.dao = newEventDAO(s)
	}
	return s.dao
}

func New(config config.DBConfig) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.config.Dsn)
	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}
