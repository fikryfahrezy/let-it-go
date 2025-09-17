package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
)

func (s *blogService) DeleteBlog(ctx context.Context, id uuid.UUID) error {
	if err := s.blogRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return repository.ErrBlogNotFound
		}
		return fmt.Errorf("%w: %w", repository.ErrFailedToDeleteBlog, err)
	}

	return nil
}