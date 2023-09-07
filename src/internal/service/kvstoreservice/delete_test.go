package kvstoreservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestDeleteWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := kvsStoreService.Delete(ctx, "key"); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestDeleteWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		deleteErr: kverror.ErrKeyNotFound,
	}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	err := kvsStoreService.Delete(context.Background(), "key")
	if err == nil {
		t.Error("error not occurred")
	}

	var kvErr *kverror.Error

	if !errors.As(err, &kvErr) {
		t.Error("error must be kverror.ErrKeyNotFound")
	}
}

func TestDelete(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]any{
			"key": "value",
		},
	}

	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	if err := kvsStoreService.Delete(context.Background(), "key"); err != nil {
		t.Error("error occurred")
	}

	_, ok := mockStorage.memoryDB["key"]
	if ok {
		t.Error("delete is not working!")
	}
}
