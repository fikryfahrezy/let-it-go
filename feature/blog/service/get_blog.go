package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
)

func (s *blogService) GetBlogByID(ctx context.Context, id uuid.UUID) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return GetBlogResponse{}, repository.ErrBlogNotFound
		}
		return GetBlogResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToGetBlog, err)
	}

	return BlogEntityToGetResponse(blog), nil
}