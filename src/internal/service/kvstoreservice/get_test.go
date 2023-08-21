package kvstoreservice_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestGetWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := kvsStoreService.Get(ctx, "key"); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestGetWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		getErr: errStorageGet,
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	if _, err := kvsStoreService.Get(context.Background(), "key"); !strings.Contains(
		err.Error(),
		"kvstoreservice.Set storage.Get",
	) {
		t.Error("error not occurred")
	}
}

func TestGet(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	if _, err := kvsStoreService.Get(context.Background(), "key"); err != nil {
		t.Error("error occurred")
	}
}
