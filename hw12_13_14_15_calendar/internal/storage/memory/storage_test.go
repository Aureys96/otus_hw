package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	t.Run("Create event", func(t *testing.T) {
		st := New()
		st.DAO().CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})
		assert.Equal(t, 1, len(st.data))
	})

	t.Run("Get event", func(t *testing.T) {
		st := New()
		st.DAO().CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})

		expected := storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		}
		actual, err := st.DAO().Get(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Update event", func(t *testing.T) {
		st := New()
		st.DAO().CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})

		err := st.DAO().Update(context.Background(), 0, storage.Event{
			Title:     "UpdatedEvent event",
			StartTime: time.Time{},
			Duration:  200,
			UserID:    1,
		})
		assert.NoError(t, err)

		expected := storage.Event{
			ID:        0,
			Title:     "UpdatedEvent event",
			StartTime: time.Time{},
			Duration:  200,
			UserID:    1,
		}
		actual, err := st.DAO().Get(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Delete event", func(t *testing.T) {
		st := New()
		st.DAO().CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})
		st.DAO().Delete(context.Background(), 0)
		assert.Equal(t, 0, len(st.data))
	})
}
