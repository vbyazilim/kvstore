package kvstoreservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestUpdateWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := kvsStoreService.Update(ctx, nil); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestUpdateWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		updateErr: kverror.ErrKeyNotFound, // raises kverror.ErrKeyNotFound
	}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	updateRequest := kvstoreservice.UpdateRequest{
		Key:   "key",
		Value: "value",
	}

	res, err := kvsStoreService.Update(context.Background(), &updateRequest)
	if res != nil {
		t.Errorf("response must be nil!")
	}

	var kvErr *kverror.Error

	if !errors.As(err, &kvErr) {
		t.Error("error must be kverror.ErrKeyNotFound")
	}
}

func TestUpdate(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]any{
			"key": "value",
		},
	}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	updateRequest := kvstoreservice.UpdateRequest{
		Key:   "key",
		Value: "vigo",
	}

	res, err := kvsStoreService.Update(context.Background(), &updateRequest)
	if err != nil {
		t.Errorf("error occurred, err: %v", err)
	}

	if res == nil {
		t.Error("result should not be nil")
	}

	if res != nil {
		val := *res

		if val.Value != "vigo" {
			t.Errorf("want: vigo, got: %s", val.Value)
		}
	}
}
