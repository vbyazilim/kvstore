package kvstorage_test

import (
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestDeleteEmpty(t *testing.T) {
	storage := kvstorage.New()

	if err := storage.Delete("key"); err == nil {
		t.Error("error not occurred")
	}
}

func TestDelete(t *testing.T) {
	key := "key"
	memoryStorage := map[string]any{
		key: "value",
	}
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	if err := storage.Delete(key); err != nil {
		t.Error("error occurred")
	}
}
