package kvstorage_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestDeleteEmpty(t *testing.T) {
	storage := kvstorage.New()

	if err := storage.Delete("key"); err == nil {
		t.Error("error not occurred")
	}
}

func TestDelete(t *testing.T) {
	key := uuid.New().String()
	memoryStorage := map[string]interface{}{
		key: "value",
	}
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	if err := storage.Delete(key); err != nil {
		t.Error("error occurred")
	}
}
