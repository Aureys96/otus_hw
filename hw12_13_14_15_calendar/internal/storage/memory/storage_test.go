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
		st.CreateEvent(context.Background(), storage.Event{
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
		st.CreateEvent(context.Background(), storage.Event{
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
		actual, err := st.Get(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Update event", func(t *testing.T) {
		st := New()
		st.CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})

		err := st.Update(context.Background(), 0, storage.Event{
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
		actual, err := st.Get(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Delete event", func(t *testing.T) {
		st := New()
		st.CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: time.Time{},
			Duration:  300,
			UserID:    1,
		})
		st.Delete(context.Background(), 0)
		assert.Equal(t, 0, len(st.data))
	})

	t.Run("booked error", func(t *testing.T) {
		st := New()
		baseDate := time.Date(2022, 8, 23, 10, 10, 10, 10, time.UTC)
		st.CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: baseDate,
			Duration:  100,
			UserID:    1,
		})
		_, err := st.CreateEvent(context.Background(), storage.Event{
			ID:        0,
			Title:     "New event",
			StartTime: baseDate.Add(10),
			Duration:  10,
			UserID:    1,
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrEventBooked)
	})
}
