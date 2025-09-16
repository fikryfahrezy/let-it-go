package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (r *blogRepository) Update(ctx context.Context, blog Blog) error {
	query := `
		UPDATE blogs
		SET title = ?, content = ?, status = ?, published_at = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	blog.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, blog.Title, blog.Content, blog.Status, blog.PublishedAt, now, blog.ID)
	if err != nil {
		slog.Error("Failed to update blog",
			slog.String("error", err.Error()),
			slog.Int("blog_id", blog.ID),
		)
		return fmt.Errorf("failed to update blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("blog not found")
	}

	slog.Info("Blog updated successfully",
		slog.Int("blog_id", blog.ID),
	)

	return nil
}