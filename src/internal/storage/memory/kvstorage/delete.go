package kvstorage

func (ms *memoryStorage) Delete(key string) error {
	if _, err := ms.Get(key); err != nil { // can not delete! key doesn't exist
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.db, key)
	return nil
}
