package kvstoreservice_test

import (
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
)

var _ kvstorage.Storer = (*mockStorage)(nil) // compile time proof

type mockStorage struct {
	deleteErr error
	getErr    error
	updateErr error
	setErr    error

	memoryDB kvstorage.MemoryDB
}

func (m *mockStorage) Delete(k string) error {
	if m.deleteErr == nil {
		delete(m.memoryDB, k)
		return nil
	}
	return m.deleteErr
}

func (m *mockStorage) Get(k string) (any, error) {
	if m.getErr == nil {
		v, ok := m.memoryDB[k]
		if !ok {
			return nil, m.getErr
		}
		return v, nil
	}
	return nil, m.getErr
}

func (m *mockStorage) List() kvstorage.MemoryDB {
	return m.memoryDB
}

func (m *mockStorage) Set(k string, v any) (any, error) {
	if m.setErr == nil {
		if _, ok := m.memoryDB[k]; ok {
			return nil, m.setErr
		}

		m.memoryDB[k] = v
		return v, nil

	}
	return nil, m.setErr
}

func (m *mockStorage) Update(k string, v any) (any, error) {
	if m.updateErr == nil {
		if _, ok := m.memoryDB[k]; !ok {
			return nil, m.updateErr
		}

		m.memoryDB[k] = v
		return v, nil
	}
	return nil, m.updateErr
}
