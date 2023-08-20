package kvstoreservice

import (
	"context"
	"fmt"
)

func (s *kvStoreService) Update(ctx context.Context, sr *UpdateRequest) (*ItemResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, err := s.storage.Update(sr.Key, sr.Value)
		if err != nil {
			return nil, fmt.Errorf("kvstoreservice.Set storage.Update err: %w", err)
		}
		return &ItemResponse{
			Key:   sr.Key,
			Value: value,
		}, nil
	}
}
