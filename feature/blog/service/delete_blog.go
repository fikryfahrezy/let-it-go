package service

import (
	"context"
	"fmt"
)

func (s *blogService) DeleteBlog(ctx context.Context, id int) error {
	if err := s.blogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}

	return nil
}