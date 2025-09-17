package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *blogService) DeleteBlog(ctx context.Context, id uuid.UUID) error {
	if err := s.blogRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}