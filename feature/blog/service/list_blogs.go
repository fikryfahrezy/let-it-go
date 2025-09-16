package service

import (
	"context"
	"fmt"
	"math"
)

func (s *blogService) ListBlogs(ctx context.Context, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error) {
	offset := (page - 1) * pageSize

	blogs, err := s.blogRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("failed to list blogs: %w", err)
	}

	totalItems, err := s.blogRepo.Count(ctx)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("failed to count blogs: %w", err)
	}

	pagination := PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(pageSize))),
	}

	return BlogEntitiesToGetResponses(blogs), pagination, nil
}