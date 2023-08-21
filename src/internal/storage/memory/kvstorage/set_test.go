package kvstorage_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/vbyazilim/kvstore/src/internal/storage"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestSet(t *testing.T) {
	key := uuid.New().String()
	memoryStorage := storage.MemoryDB(map[string]interface{}{})
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	storage.Set(key, "value")

	if _, ok := memoryStorage[key]; !ok {
		t.Error("value not equal")
	}
}
