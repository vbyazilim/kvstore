package kvstorage

import (
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

func (ms *memoryStorage) Get(key string) (any, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	value, ok := ms.db[key]
	if !ok {
		return nil, fmt.Errorf("%w, %q not exists", kverror.ErrKeyNotFound, key)
	}
	return value, nil
}
