package service

import (
	"context"
	"fmt"
	"math"
)

func (s *blogService) GetBlogsByStatus(ctx context.Context, status string, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error) {
	offset := (page - 1) * pageSize

	blogs, err := s.blogRepo.GetByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("failed to get blogs by status: %w", err)
	}

	totalItems, err := s.blogRepo.CountByStatus(ctx, status)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("failed to count blogs by status: %w", err)
	}

	pagination := PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(pageSize))),
	}

	return BlogEntitiesToGetResponses(blogs), pagination, nil
}