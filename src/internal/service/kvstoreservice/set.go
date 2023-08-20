package kvstoreservice

import (
	"context"
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/service"
)

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
