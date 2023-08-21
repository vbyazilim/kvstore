package kvstoreservice_test

import (
	"errors"

	"github.com/vbyazilim/kvstore/src/internal/storage"
)

var (
	errStorageDelete = errors.New("storage delete error")
	errStorageGet    = errors.New("storage get error")
	errStorageSet    = errors.New("storage set error")
)

type mockStorage struct {
	deleteErr error
	getErr    error
	setErr    error
	updateErr error
	memoryDB  storage.MemoryDB
}

func (m *mockStorage) Delete(_ string) error {
	return m.deleteErr
}

func (m *mockStorage) Get(_ string) (interface{}, error) {
	return nil, m.getErr
}

func (m *mockStorage) List() storage.MemoryDB {
	return m.memoryDB
}

func (m *mockStorage) Set(_ string, _ interface{}) (interface{}, error) {
	return nil, m.setErr
}

func (m *mockStorage) Update(_ string, _ interface{}) (interface{}, error) {
	return nil, m.updateErr
}
