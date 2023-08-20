package kvstoreservice

import (
	"context"
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/service"
)

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
