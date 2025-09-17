package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
)

func (s *blogService) PublishBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return GetBlogResponse{}, repository.ErrBlogNotFound
		}
		return GetBlogResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToGetBlog, err)
	}

	blog.Status = repository.StatusPublished
	now := time.Now()
	blog.PublishedAt = &now

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return GetBlogResponse{}, repository.ErrBlogNotFound
		}
		return GetBlogResponse{}, fmt.Errorf("%w: %w", ErrFailedToPublishBlog, err)
	}

	return BlogEntityToGetResponse(blog), nil
}