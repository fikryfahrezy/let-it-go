package service

import (
	"context"

	"github.com/google/uuid"
)

type BlogService interface {
	CreateBlog(ctx context.Context, req CreateBlogRequest) (GetBlogResponse, error)
	GetBlogByID(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
	GetBlogsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error)
	GetBlogsByStatus(ctx context.Context, status string, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error)
	UpdateBlog(ctx context.Context, id uuid.UUID, req UpdateBlogRequest) (GetBlogResponse, error)
	DeleteBlog(ctx context.Context, id uuid.UUID) error
	ListBlogs(ctx context.Context, page, pageSize int) ([]GetBlogResponse, PaginationInfo, error)
	PublishBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
	ArchiveBlog(ctx context.Context, id uuid.UUID) (GetBlogResponse, error)
}

type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}