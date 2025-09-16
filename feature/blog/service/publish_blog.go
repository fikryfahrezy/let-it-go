package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func (s *blogService) PublishBlog(ctx context.Context, id int) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to get blog: %w", err)
	}

	blog.Status = repository.StatusPublished
	now := time.Now()
	blog.PublishedAt = &now

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to publish blog: %w", err)
	}

	return BlogEntityToGetResponse(blog), nil
}