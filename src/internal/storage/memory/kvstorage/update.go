package kvstorage

func (ms *memoryStorage) Update(key string, value any) (any, error) {
	if _, err := ms.Get(key); err != nil { // can not update! key doesn't exist
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.db[key] = value
	return value, nil
}
