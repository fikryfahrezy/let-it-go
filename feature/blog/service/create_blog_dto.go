package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type CreateBlogRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID uuid.UUID `json:"author_id"`
	Status   string `json:"status"`
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