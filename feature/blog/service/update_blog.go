package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *blogService) UpdateBlog(ctx context.Context, id uuid.UUID, req UpdateBlogRequest) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, err
	}

	req.ApplyToEntity(blog)

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return GetBlogResponse{}, err
	}

	return BlogEntityToGetResponse(blog), nil
}
