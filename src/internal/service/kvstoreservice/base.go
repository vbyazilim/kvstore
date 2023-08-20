package kvstoreservice

import (
	"context"

	"github.com/vbyazilim/kvstore/src/internal/storage"
)

var _ KVStoreService = (*kvStoreService)(nil) // compile time proof

// KVStoreService defines service behaviours.
type KVStoreService interface {
	Set(context.Context, *SetRequest) (*ItemResponse, error)
	Get(context.Context, string) (*ItemResponse, error)
	Update(context.Context, *UpdateRequest) (*ItemResponse, error)
	Delete(context.Context, string) error
	List(context.Context) (*ListResponse, error)
}

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
func New(options ...ServiceOption) KVStoreService {
	kvs := &kvStoreService{}

	for _, o := range options {
		o(kvs)
	}

	return kvs
}
