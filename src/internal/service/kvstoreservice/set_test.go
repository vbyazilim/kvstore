package kvstoreservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func TestSetWithCancel(t *testing.T) {
	mockStorage := &mockStorage{}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := kvsStoreService.Set(ctx, nil); !errors.Is(err, ctx.Err()) {
		t.Error("error not occurred")
	}
}

func TestSetWithStorageError(t *testing.T) {
	mockStorage := &mockStorage{
		setErr: kverror.ErrKeyExists,
	}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	serviceRequest := kvstoreservice.SetRequest{
		Key:   "vigo",
		Value: "lego",
	}

	res, err := kvsStoreService.Set(context.Background(), &serviceRequest)

	if res != nil {
		t.Errorf("response must be nil!")
	}

	var kvErr *kverror.Error

	if !errors.As(err, &kvErr) {
		t.Error("error must be kverror.ErrKeyExists")
	}
}

func TestSet(t *testing.T) {
	mockStorage := &mockStorage{
		memoryDB: map[string]any{},
	}
	kvsStoreService := kvstoreservice.New(
		kvstoreservice.WithStorage(mockStorage),
	)

	setRequest := kvstoreservice.SetRequest{
		Key:   "username",
		Value: "vigo",
	}

	res, err := kvsStoreService.Set(context.Background(), &setRequest)
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
