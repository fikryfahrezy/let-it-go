package service

//counterfeiter:generate -o servicefakes/fake_blog_service.go . BlogService

import (
	"context"

	"github.com/google/uuid"
)

type BlogService interface {
	CreateBlog(ctx context.Context, req CreateBlogRequest) (GetBlogResponse, error)
	GetBlogByID(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
	GetBlogsByAuthor(ctx context.Context, authorID uuid.UUID, req GetBlogsByAuthorRequest) ([]GetBlogResponse, int, error)
	GetBlogsByStatus(ctx context.Context, status string, req GetBlogsByStatusRequest) ([]GetBlogResponse, int, error)
	UpdateBlog(ctx context.Context, id uuid.UUID, req UpdateBlogRequest) (GetBlogResponse, error)
	DeleteBlog(ctx context.Context, id uuid.UUID) error
	ListBlogs(ctx context.Context, req ListBlogsRequest) ([]GetBlogResponse, int, error)
	PublishBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
	ArchiveBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
}
