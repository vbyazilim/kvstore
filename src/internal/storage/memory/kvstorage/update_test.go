package kvstorage_test

import (
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestUpdateEmpty(t *testing.T) {
	storage := kvstorage.New()

	if _, err := storage.Update("key", "value"); err == nil {
		t.Error("error not occurred")
	}
}

func TestUpdate(t *testing.T) {
	key := "key"
	memoryStorage := map[string]any{
		key: "value",
	}
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	value, err := storage.Update(key, "value2")
	if err != nil {
		t.Error("error occurred")
	}

	if value != "value2" {
		t.Error("value not equal")
	}
}
