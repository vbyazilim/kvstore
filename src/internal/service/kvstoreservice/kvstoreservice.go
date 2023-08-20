package kvstoreservice

import (
	"context"
	"fmt"

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

func (s *kvStoreService) Set(ctx context.Context, sr *service.SetRequest) (*service.ItemResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, err := s.storage.Set(sr.Key, sr.Value)
		if err != nil {
			return nil, fmt.Errorf("kvstoreservice.Set storage.Set err: %w", err)
		}
		return &service.ItemResponse{
			Key:   sr.Key,
			Value: value,
		}, nil
	}
}

func (s *kvStoreService) Get(ctx context.Context, key string) (*service.ItemResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, err := s.storage.Get(key)
		if err != nil {
			return nil, fmt.Errorf("kvstoreservice.Set storage.Get err: %w", err)
		}
		return &service.ItemResponse{
			Key:   key,
			Value: value,
		}, nil
	}
}

func (s *kvStoreService) Update(ctx context.Context, sr *service.UpdateRequest) (*service.ItemResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, err := s.storage.Update(sr.Key, sr.Value)
		if err != nil {
			return nil, fmt.Errorf("kvstoreservice.Set storage.Update err: %w", err)
		}
		return &service.ItemResponse{
			Key:   sr.Key,
			Value: value,
		}, nil
	}
}

func (s *kvStoreService) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.storage.Delete(key); err != nil {
			return fmt.Errorf("kvstoreservice.Set storage.Delete err: %w", err)
		}
		return nil
	}
}

func (s *kvStoreService) List(ctx context.Context) (*service.ListResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		items := s.storage.List()
		response := make(service.ListResponse, len(items))
		for k, v := range items {
			response = append(response, service.ItemResponse{
				Key:   k,
				Value: v,
			})
		}
		return &response, nil
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
