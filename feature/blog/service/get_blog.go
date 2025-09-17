package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *blogService) GetBlogByID(ctx context.Context, id uuid.UUID) (GetBlogResponse, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return GetBlogResponse{}, err
	}

	return BlogEntityToGetResponse(blog), nil
}
