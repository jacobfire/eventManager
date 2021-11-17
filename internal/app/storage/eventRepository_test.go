package storage_test

import (
	"calendar/internal/app/model"
	"calendar/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventRepository_Load(t *testing.T) {
	var createdId int
	t.Run("Check Creation", func (t *testing.T) {
		event := &model.Event{
			Title: "test element",
			Description: "test description",
			Time: "2021-10-11 2:45 PM",
			Timezone: "Europe/Kiev",
			Duration: 24,
			Notes: []string{"test", "the", "event"},
		}

		s, _ := storage.TestStorage(t, databaseURL)
		e, err := s.Event().Create(event)
		//defer teardown("event")

		createdId = e.ID
		assert.NoError(t, err)
		assert.NotNil(t, e)
	})

	t.Run("Check Removing", func (t *testing.T) {
		s, _ := storage.TestStorage(t, databaseURL)
		b, err := s.Event().Delete(createdId)

		assert.NoError(t, err)
		assert.True(t, b)
		//assert.True(t, b)
	})
}

