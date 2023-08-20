package kvstorage

import (
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

func (ms *memoryStorage) Set(key string, value any) (any, error) {
	val, err := ms.Get(key)
	if err == nil { // this means, key already exists!
		dataMsg := "can not set '" + key + "'"
		vals, ok := val.(string)
		if ok {
			dataMsg = dataMsg + ", value '" + vals + "'"
		}
		return nil, fmt.Errorf("%w", kverror.ErrKeyExists.AddData(dataMsg))
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.db[key] = value
	return value, nil
}
