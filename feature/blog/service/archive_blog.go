package service

import (
	"context"
	"fmt"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func (s *blogService) ArchiveBlog(ctx context.Context, id int) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to get blog: %w", err)
	}

	blog.Status = repository.StatusArchived
	blog.PublishedAt = nil

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to archive blog: %w", err)
	}

	return BlogEntityToGetResponse(blog), nil
}