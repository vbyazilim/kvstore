package kvstorage_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/vbyazilim/kvstore/src/internal/storage"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

func TestList(t *testing.T) {
	key := uuid.New().String()
	memoryStorage := storage.MemoryDB(map[string]interface{}{
		key: "value",
	})
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(memoryStorage),
	)

	value := storage.List()

	if !reflect.DeepEqual(value, memoryStorage) {
		t.Error("value not equal")
	}
}
