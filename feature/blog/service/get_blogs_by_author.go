package service

import (
	"context"
	"fmt"
	"math"
)

func (s *blogService) GetBlogsByAuthor(ctx context.Context, authorID, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error) {
	offset := (page - 1) * pageSize

	blogs, err := s.blogRepo.GetByAuthorID(ctx, authorID, pageSize, offset)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("failed to get blogs by author: %w", err)
	}

	// Get total count for pagination (you might want to add a CountByAuthorID method)
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