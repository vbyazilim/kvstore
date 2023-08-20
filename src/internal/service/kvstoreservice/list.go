package kvstoreservice

import (
	"context"

	"github.com/vbyazilim/kvstore/src/internal/service"
)

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
