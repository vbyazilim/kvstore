package kvstorage

func (ms *memoryStorage) Set(key string, value any) (any, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.db[key] = value
	return value, nil
}
