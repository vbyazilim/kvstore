package kvstoreservice_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestSetWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := kvsStoreService.Set(ctx, nil)
	if !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestSetWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		setErr: errStorageSet,
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	setRequest := kvstoreservice.SetRequest{
		Key:   "key",
		Value: "value",
	}
	_, err := kvsStoreService.Set(context.Background(), &setRequest)
	if !strings.Contains(err.Error(), "kvstoreservice.Set storage.Set") {
		t.Error("error not occurred")
	}
}

func TestSet(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]interface{}{
			"key": "value",
		},
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	setRequest := kvstoreservice.SetRequest{
		Key:   "key",
		Value: "value",
	}
	_, err := kvsStoreService.Set(context.Background(), &setRequest)
	if err != nil {
		t.Error("error occurred")
	}
}
