package kvstorage

import (
	"fmt"
	"sync"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
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

func (ms *memoryStorage) Set(key string, value any) (any, error) {
	val, err := ms.Get(key)
	if err == nil { // this means, key already exists!
		return nil, fmt.Errorf("%w can not set %q, value exists: %v", kverror.ErrKeyExists, key, val)
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.db[key] = value
	return value, nil
}

func (ms *memoryStorage) Update(key string, value any) (any, error) {
	if _, err := ms.Get(key); err != nil { // can not update! key doesn't exist
		return nil, err
	}
	return ms.Set(key, value)
}

func (ms *memoryStorage) Get(key string) (any, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	value, ok := ms.db[key]
	if !ok {
		return nil, fmt.Errorf("%w, %q not exists", kverror.ErrKeyNotFound, key)
	}
	return value, nil
}

func (ms *memoryStorage) Delete(key string) error {
	if _, err := ms.Get(key); err != nil { // can not delete! key doesn't exist
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.db, key)
	return nil
}

func (ms *memoryStorage) List() storage.MemoryDB {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.db
}

// New instantiates new storage instance.
func New(options ...StorageOption) storage.Storer {
	ms := &memoryStorage{}

	for _, o := range options {
		o(ms)
	}

	return ms
}
