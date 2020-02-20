package adapters

import (
	"context"
	"sync"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/google/uuid"
)

// UUIDTable is a mapping table between authenticational credentials (login,password,ip) and UUID
type UUIDTable struct {
	table map[string]uuid.UUID
	mx    sync.RWMutex
}

// NewUUIDTable returns a new instance of the UUIDTable
func NewUUIDTable() *UUIDTable {
	return &UUIDTable{
		table: map[string]uuid.UUID{},
		mx:    sync.RWMutex{},
	}
}

// AddToTable adds a new credential to the table, in case if it's already exists returns it's UUID & corresponding error
func (u *UUIDTable) AddToTable(ctx context.Context, value string) uuid.UUID {
	u.mx.Lock()
	defer u.mx.Unlock()

	if u.isPresentInTable(ctx, value) {
		return u.table[value]
	}

	u.table[value] = uuid.New()

	return u.table[value]
}

// DeleteFromTable deletes a credential to the table, in case if it's already exists returns it's UUID & corresponding error
func (u *UUIDTable) DeleteFromTable(ctx context.Context, value string) error {
	u.mx.RLock()
	defer u.mx.RUnlock()

	if u.isPresentInTable(ctx, value) {
		delete(u.table, value)
		return nil
	}

	return errors.ErrNoMappingFound
}

func (u *UUIDTable) isPresentInTable(ctx context.Context, value string) bool {
	if _, ok := u.table[value]; ok {
		return true
	}

	return false
}
