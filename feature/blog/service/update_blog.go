package service

import (
	"context"
	"fmt"
)

func (s *blogService) UpdateBlog(ctx context.Context, id int, req UpdateBlogRequest) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to get blog: %w", err)
	}

	req.ApplyToEntity(blog)

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to update blog: %w", err)
	}

	return BlogEntityToGetResponse(blog), nil
}