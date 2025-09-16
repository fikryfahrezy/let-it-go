package service

import (
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type GetBlogResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	AuthorID    int        `json:"author_id"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func BlogEntityToGetResponse(blog repository.Blog) GetBlogResponse {
	return GetBlogResponse{
		ID:          blog.ID,
		Title:       blog.Title,
		Content:     blog.Content,
		AuthorID:    blog.AuthorID,
		Status:      blog.Status,
		PublishedAt: blog.PublishedAt,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
	}
}

func BlogEntitiesToGetResponses(blogs []repository.Blog) []GetBlogResponse {
	responses := make([]GetBlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = BlogEntityToGetResponse(blog)
	}
	return responses
}
