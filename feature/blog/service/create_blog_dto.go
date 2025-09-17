package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type CreateBlogRequest struct {
	Title    string    `json:"title" validate:"required,min=3,max=200"`
	Content  string    `json:"content" validate:"required,min=10"`
	AuthorID uuid.UUID `json:"author_id" validate:"required"`
	Status   string    `json:"status" validate:"required,oneof=draft published archived"`
}

func (req CreateBlogRequest) ToEntity() repository.Blog {
	blog := repository.Blog{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: req.AuthorID,
		Status:   req.Status,
	}

	// Set published_at if status is published
	if req.Status == repository.StatusPublished {
		now := time.Now()
		blog.PublishedAt = &now
	}

	return blog
}