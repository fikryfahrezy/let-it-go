package service

import (
	"context"
	"fmt"
	"math"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func (s *blogService) ListBlogs(ctx context.Context, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error) {
	offset := (page - 1) * pageSize

	blogs, err := s.blogRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("%w: %w", repository.ErrFailedToListBlogs, err)
	}

	totalItems, err := s.blogRepo.Count(ctx)
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("%w: %w", repository.ErrFailedToCountBlogs, err)
	}

	pagination := PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(pageSize))),
	}

	return BlogEntitiesToGetResponses(blogs), pagination, nil
}