package kvstoreservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestListWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := kvsStoreService.List(ctx); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestList(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]any{
			"key": "value",
		},
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	if _, err := kvsStoreService.List(context.Background()); err != nil {
		t.Error("error occurred")
	}
}
