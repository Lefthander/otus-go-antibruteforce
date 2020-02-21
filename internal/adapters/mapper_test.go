package adapters

import (
	"context"
	"testing"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	var id uuid.UUID

	tb := NewUUIDTable()

	ctx := context.Background()

	t.Run("Check empty mapper table", func(t *testing.T) {
		assert.Equal(t, 0, len(tb.table))
	})

	t.Run("Add new record to the Mapper", func(t *testing.T) {
		id = tb.AddToTable(ctx, "Value1")
		assert.Equal(t, 1, len(tb.table))
	})

	t.Run("Add the same value to the Mapper", func(t *testing.T) {
		newid := tb.AddToTable(ctx, "Value1")
		assert.Equal(t, id, newid)
		assert.Equal(t, 1, len(tb.table))
	})

	t.Run("Delete from the Mapper", func(t *testing.T) {
		err := tb.DeleteFromTable(ctx, "Value1")
		assert.Equal(t, nil, err)
		assert.Equal(t, 0, len(tb.table))
	})

	t.Run("Delete from the empty Mapper", func(t *testing.T) {
		err := tb.DeleteFromTable(ctx, "Value1")
		assert.Equal(t, errors.ErrNoMappingFound, err)
	})

	t.Run("Clear the Mapper", func(t *testing.T) {
		_ = tb.AddToTable(ctx, "Value1")
		_ = tb.AddToTable(ctx, "Value2")
		_ = tb.AddToTable(ctx, "Value3")
		assert.Equal(t, 3, len(tb.table))
		tb.Clear()
		assert.Equal(t, 0, len(tb.table))
	})
}
