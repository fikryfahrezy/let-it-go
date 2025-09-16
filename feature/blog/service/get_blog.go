package service

import (
	"context"
	"fmt"
)

func (s *blogService) GetBlogByID(ctx context.Context, id int) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, fmt.Errorf("failed to get blog: %w", err)
	}

	return BlogEntityToGetResponse(blog), nil
}