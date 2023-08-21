package kvstoreservice_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestDeleteWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := kvsStoreService.Delete(ctx, "key"); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestDeleteWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		deleteErr: errStorageDelete,
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	if err := kvsStoreService.Delete(context.Background(), "key"); !strings.Contains(
		err.Error(),
		"kvstoreservice.Set storage.Delete",
	) {
		t.Error("error not occurred")
	}
}

func TestDelete(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	if err := kvsStoreService.Delete(context.Background(), "key"); err != nil {
		t.Error("error occurred")
	}
}
