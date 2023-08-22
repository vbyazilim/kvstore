package kvstorage_test

import (
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestSet(t *testing.T) {
	key := "key"
	memoryStorage := kvstorage.MemoryDB(map[string]any{})
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	storage.Set(key, "value")

	if _, ok := memoryStorage[key]; !ok {
		t.Error("value not equal")
	}
}
