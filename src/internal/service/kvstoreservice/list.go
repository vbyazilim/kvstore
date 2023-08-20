package kvstoreservice

import (
	"context"
)

func (s *kvStoreService) List(ctx context.Context) (*ListResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		items := s.storage.List()
		response := make(ListResponse, len(items))
		for k, v := range items {
			response = append(response, ItemResponse{
				Key:   k,
				Value: v,
			})
		}
		return &response, nil
	}
}
