package adapters

import (
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

// AddToTable adds a new credential to the table, in case if it's already exists
// returns it's UUID
func (u *UUIDTable) AddToTable(value string) uuid.UUID {
	if u.isPresentInTable(value) {
		return u.table[value]
	}

	u.mx.RLock()
	u.table[value] = uuid.New()
	u.mx.RUnlock()

	return u.table[value]
}

// DeleteFromTable deletes a credential from the table, in case if it's not found
// returns corresponding error
func (u *UUIDTable) DeleteFromTable(value string) error {
	if u.isPresentInTable(value) {
		u.mx.RLock()
		delete(u.table, value)
		u.mx.RUnlock()

		return nil
	}

	return errors.ErrNoMappingFound
}

func (u *UUIDTable) isPresentInTable(value string) bool {
	u.mx.Lock()
	defer u.mx.Unlock()

	if _, ok := u.table[value]; ok {
		return true
	}

	return false
}

// Clear delete all records from the table
func (u *UUIDTable) Clear() {
	u.mx.Lock()

	for k := range u.table {
		delete(u.table, k)
	}

	u.mx.Unlock()
}
