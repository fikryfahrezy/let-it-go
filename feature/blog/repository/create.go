package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (r *blogRepository) Create(ctx context.Context, blog Blog) error {
	query := `
		INSERT INTO blogs (title, content, author_id, status, published_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, blog.Title, blog.Content, blog.AuthorID, blog.Status, blog.PublishedAt, now, now)
	if err != nil {
		slog.Error("Failed to create blog",
			slog.String("error", err.Error()),
			slog.String("title", blog.Title),
		)
		return fmt.Errorf("failed to create blog: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert id",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	blog.ID = int(id)

	slog.Info("Blog created successfully",
		slog.Int("blog_id", blog.ID),
		slog.String("title", blog.Title),
	)

	return nil
}