package service

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func (s *blogService) CreateBlog(ctx context.Context, req CreateBlogRequest) (GetBlogResponse, error) {
	if req.Status == "" {
		req.Status = repository.StatusDraft
	}

	blog := req.ToEntity()

	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return GetBlogResponse{}, err
	}

	return BlogEntityToGetResponse(blog), nil
}