package kvstoreservice

import (
	"github.com/vbyazilim/kvstore/src/internal/service"
	"github.com/vbyazilim/kvstore/src/internal/storage"
)

var _ service.Servicer = (*kvStoreService)(nil) // compile time proof

type kvStoreService struct {
	storage storage.Storer
}

// ServiceOption represents service option type.
type ServiceOption func(*kvStoreService)

// WithStorage sets storage option.
func WithStorage(strg storage.Storer) ServiceOption {
	return func(s *kvStoreService) {
		s.storage = strg
	}
}

// New instantiates new service instance.
func New(options ...ServiceOption) service.Servicer {
	kvs := &kvStoreService{}

	for _, o := range options {
		o(kvs)
	}

	return kvs
}
