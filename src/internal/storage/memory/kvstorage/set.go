package kvstorage

import (
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

func (ms *memoryStorage) Set(key string, value any) (any, error) {
	if _, err := ms.Get(key); err == nil {
		return nil, fmt.Errorf("%w", kverror.ErrKeyExists.AddData("'"+key+"' already exist"))
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.db[key] = value
	return value, nil
}
