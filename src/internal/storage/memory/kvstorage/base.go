package kvstorage

import (
	"sync"
)

var _ Storer = (*memoryStorage)(nil) // compile time proof

// MemoryDB is a type alias for in memory-db type.
type MemoryDB map[string]any

// Storer defines storage behaviours.
type Storer interface {
	Set(key string, value any) any
	Get(key string) (any, error)
	Update(key string, value any) (any, error)
	Delete(key string) error
	List() MemoryDB
}

type memoryStorage struct {
	mu sync.RWMutex // guarding db only
	db MemoryDB
}

// StorageOption represents storage option type.
type StorageOption func(*memoryStorage)

// WithMemoryDB sets db option.
func WithMemoryDB(db MemoryDB) StorageOption {
	return func(s *memoryStorage) {
		s.db = db
	}
}

// New instantiates new storage instance.
func New(options ...StorageOption) Storer {
	ms := &memoryStorage{}

	for _, o := range options {
		o(ms)
	}

	return ms
}
