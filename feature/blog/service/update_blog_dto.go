package service

import (
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type UpdateBlogRequest struct {
	Title   string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Content string `json:"content,omitempty" validate:"omitempty,min=10"`
	Status  string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

func (req UpdateBlogRequest) ApplyToEntity(blog repository.Blog) {
	if req.Title != "" {
		blog.Title = req.Title
	}
	if req.Content != "" {
		blog.Content = req.Content
	}
	if req.Status != "" {
		// Handle status changes
		oldStatus := blog.Status
		blog.Status = req.Status

		// Set published_at when changing to published
		if req.Status == repository.StatusPublished && oldStatus != repository.StatusPublished {
			now := time.Now()
			blog.PublishedAt = &now
		}

		// Clear published_at when changing from published to other status
		if req.Status != repository.StatusPublished && oldStatus == repository.StatusPublished {
			blog.PublishedAt = nil
		}
	}
}
