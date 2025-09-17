package service

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
)

func (s *blogService) ArchiveBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, err
	}

	blog.Status = repository.StatusArchived
	blog.PublishedAt = nil

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, err
	}

	return BlogEntityToGetResponse(blog), nil
}
