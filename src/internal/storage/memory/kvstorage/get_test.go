package kvstorage_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestGetEmpty(t *testing.T) {
	storage := kvstorage.New()

	if _, err := storage.Get("key"); err == nil {
		t.Error("error not occurred")
	}
}

func TestGet(t *testing.T) {
	key := uuid.New().String()
	memoryStorage := map[string]interface{}{
		key: "value",
	}
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	value, err := storage.Get(key)
	if err != nil {
		t.Error("error occurred")
	}

	if value != "value" {
		t.Error("value not equal")
	}
}
