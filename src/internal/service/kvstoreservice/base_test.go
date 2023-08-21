package kvstoreservice_test

import "github.com/vbyazilim/kvstore/src/internal/storage"

type mockStorage struct {
	deleteErr error
	getErr    error
	setErr    error
	updateErr error
	memoryDB  storage.MemoryDB
}

func (m *mockStorage) Delete(key string) error {
	return m.deleteErr
}

func (m *mockStorage) Get(key string) (interface{}, error) {
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
