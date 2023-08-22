package kvstoreservice

import (
	"context"
	"fmt"
)

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
