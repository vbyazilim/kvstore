package kvstorage

func (ms *memoryStorage) List() MemoryDB {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.db
}
