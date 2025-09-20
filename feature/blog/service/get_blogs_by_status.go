package service

import (
	"context"
)

func (s *blogService) GetBlogsByStatus(ctx context.Context, status string, req GetBlogsByStatusRequest) ([]GetBlogResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	blogs, err := s.blogRepo.GetByStatus(ctx, status, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.blogRepo.CountByStatus(ctx, status)
	if err != nil {
		return nil, 0, err
	}

	return BlogEntitiesToGetResponses(blogs), totalItems, nil
}
