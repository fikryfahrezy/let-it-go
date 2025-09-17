package service

import (
	"context"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
)

func (s *blogService) PublishBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, err
	}

	blog.Status = repository.StatusPublished
	now := time.Now()
	blog.PublishedAt = &now

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, err
	}

	return BlogEntityToGetResponse(blog), nil
}