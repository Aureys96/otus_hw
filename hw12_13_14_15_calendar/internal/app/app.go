package app

import (
	"context"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	storage storage.IStorage
	logger  *zap.Logger
}

func New(logger *zap.Logger, storage storage.IStorage) *App {
	return &App{storage, logger}
}

func (a App) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := a.storage.CreateEvent(ctx, event)
	return err
}

func (a App) Get(ctx context.Context, id int) (storage.Event, error) {
	return a.storage.Get(ctx, id)
}

func (a App) Update(ctx context.Context, id int, event storage.Event) error {
	return a.storage.Update(ctx, id, event)
}

func (a App) Delete(ctx context.Context, id int) {
	a.storage.Delete(ctx, id)
}

func (a App) ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	return a.storage.ListEvents(ctx, start, end)
}
