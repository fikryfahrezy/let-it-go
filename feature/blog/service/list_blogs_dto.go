package service

import "github.com/fikryfahrezy/let-it-go/pkg/http_server"

// ListBlogsRequest represents the request for listing blogs with pagination
type ListBlogsRequest struct {
	http_server.PaginationRequest
}

// GetBlogsByAuthorRequest represents the request for getting blogs by author with pagination
type GetBlogsByAuthorRequest struct {
	http_server.PaginationRequest
}

// GetBlogsByStatusRequest represents the request for getting blogs by status with pagination
type GetBlogsByStatusRequest struct {
	http_server.PaginationRequest
}
