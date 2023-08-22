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
		return nil, fmt.Errorf("%w", kverror.ErrKeyNotFound.AddData("'"+key+"' does not exist"))
	}
	return value, nil
}
