package kvstoreservice

import (
	"context"
	"fmt"

	"github.com/vbyazilim/kvstore/src/internal/service"
)

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
