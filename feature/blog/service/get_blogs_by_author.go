package service

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
)

func (s *blogService) GetBlogsByAuthor(ctx context.Context, authorID uuid.UUID, req http_server.PaginationRequest) ([]GetBlogResponse, int, error) {
	offset := (req.Page - 1) * req.PageSize

	blogs, err := s.blogRepo.GetByAuthorID(ctx, authorID, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.blogRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return BlogEntitiesToGetResponses(blogs), totalItems, nil
}