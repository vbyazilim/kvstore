package kvstorage

import "github.com/vbyazilim/kvstore/src/internal/storage"

func (ms *memoryStorage) List() storage.MemoryDB {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.db
}
