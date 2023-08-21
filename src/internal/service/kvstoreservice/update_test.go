package kvstoreservice_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestUpdateWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := kvsStoreService.Update(ctx, nil)
	if !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestUpdateWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		updateErr: errStorageUpdate,
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	updateRequest := kvstoreservice.UpdateRequest{
		Key:   "key",
		Value: "value",
	}
	_, err := kvsStoreService.Update(context.Background(), &updateRequest)
	if !strings.Contains(err.Error(), "kvstoreservice.Set storage.Update") {
		t.Error("error not occurred")
	}
}

func TestUpdate(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]interface{}{
			"key": "value",
		},
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	updateRequest := kvstoreservice.UpdateRequest{
		Key:   "key",
		Value: "value",
	}
	_, err := kvsStoreService.Update(context.Background(), &updateRequest)
	if err != nil {
		t.Error("error occurred")
	}
}
