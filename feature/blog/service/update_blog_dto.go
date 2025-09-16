package service

import (
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type UpdateBlogRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
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
