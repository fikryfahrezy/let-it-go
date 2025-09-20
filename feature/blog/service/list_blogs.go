package service

import (
	"context"
)

func (s *blogService) ListBlogs(ctx context.Context, req ListBlogsRequest) ([]GetBlogResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	blogs, err := s.blogRepo.List(ctx, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := s.blogRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return BlogEntitiesToGetResponses(blogs), totalItems, nil
}
