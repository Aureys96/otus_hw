package app

import (
	"context"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
	"time"
)

type App struct {
	storage storage.IStorage
	logger  *zap.Logger
}

func (a App) CreateEvent(ctx context.Context, event storage.Event) error {
	//TODO implement me
	panic("implement me")
}

func (a App) Get(ctx context.Context, id int) (storage.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a App) Update(ctx context.Context, id int, event storage.Event) error {
	//TODO implement me
	panic("implement me")
}

func (a App) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (a App) ListEvents(ctx context.Context, date time.Time) []storage.Event {
	//TODO implement me
	panic("implement me")
}

func New(logger *zap.Logger, storage storage.IStorage) *App {
	return &App{storage, logger}
}
