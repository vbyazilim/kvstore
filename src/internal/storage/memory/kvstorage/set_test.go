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

	val, err := storage.Set(key, "value")
	if err != nil {
		t.Errorf("want: value, got: %v, err: %v", val, err)
	}

	if _, err := storage.Set(key, "xxx"); err == nil {
		t.Error("error not occurred")
	}
}
