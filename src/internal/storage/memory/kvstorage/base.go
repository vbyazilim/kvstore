package kvstorage

import (
	"sync"

	"github.com/vbyazilim/kvstore/src/internal/storage"
)

var _ storage.Storer = (*memoryStorage)(nil) // compile time proof

type memoryStorage struct {
	mu sync.RWMutex // guarding db only
	db storage.MemoryDB
}

// StorageOption represents storage option type.
type StorageOption func(*memoryStorage)

// WithMemoryDB sets db option.
func WithMemoryDB(db storage.MemoryDB) StorageOption {
	return func(s *memoryStorage) {
		s.db = db
	}
}

// New instantiates new storage instance.
func New(options ...StorageOption) storage.Storer {
	ms := &memoryStorage{}

	for _, o := range options {
		o(ms)
	}

	return ms
}
