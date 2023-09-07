package kvstoreservice

import (
	"context"
	"fmt"
)

func (s *kvStoreService) Set(ctx context.Context, sr *SetRequest) (*ItemResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if _, err := s.storage.Set(sr.Key, sr.Value); err != nil {
			return nil, fmt.Errorf("kvstoreservice.Set storage.Set err: %w", err)
		}

		return &ItemResponse{
			Key:   sr.Key,
			Value: sr.Value,
		}, nil
	}
}
