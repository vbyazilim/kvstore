package kvstoreservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
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
		getErr: kverror.ErrKeyNotFound, // get raises ErrKeyNotFound
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	res, err := kvsStoreService.Get(context.Background(), "key")
	if err == nil {
		t.Error("error not occurred")
	}

	if res != nil {
		t.Errorf("response must be nil!")
	}

	var kvErr *kverror.Error

	if !errors.As(err, &kvErr) {
		t.Error("error must be kverror.ErrKeyNotFound")
	}
}

func TestGet(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]any{
			"key": "value",
		},
	}
	kvsStoreService := kvstoreservice.New(kvstoreservice.WithStorage(mockStorage))

	res, err := kvsStoreService.Get(context.Background(), "key")
	if err != nil {
		t.Error("error occurred")
	}

	if res == nil {
		t.Error("result should not be nil")
	}

	if res != nil {
		val := *res
		if val.Value != "value" {
			t.Errorf("want: value, got: %s", val.Value)
		}
	}
}
