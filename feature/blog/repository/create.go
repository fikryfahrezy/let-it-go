package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (r *blogRepository) Create(ctx context.Context, blog Blog) error {
	query := `
		INSERT INTO blogs (id, title, content, author_id, status, published_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now

	// Generate UUIDv7 for the blog ID
	blog.ID = uuid.Must(uuid.NewV7())

	_, err := r.db.ExecContext(ctx, query, blog.ID, blog.Title, blog.Content, blog.AuthorID, blog.Status, blog.PublishedAt, now, now)
	if err != nil {
		r.log.Error("Failed to create blog",
			slog.String("error", err.Error()),
			slog.String("title", blog.Title),
		)
		return fmt.Errorf("%w: %w", ErrFailedToCreateBlog, err)
	}

	// No need to get last insert ID since we're using UUIDs

	r.log.Info("Blog created successfully",
		slog.String("blog_id", blog.ID.String()),
		slog.String("title", blog.Title),
	)

	return nil
}
