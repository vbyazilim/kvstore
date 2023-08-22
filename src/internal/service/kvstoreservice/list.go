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

		var i int
		for k, v := range items {
			response[i] = ItemResponse{
				Key:   k,
				Value: v,
			}
			i++
		}
		return &response, nil
	}
}
