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
		r.log.Error("Failed to update blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", blog.ID.String()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToUpdateBlog, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToGetRowsAffected, err)
	}

	if rowsAffected == 0 {
		return ErrBlogNotFound
	}

	r.log.Info("Blog updated successfully",
		slog.String("blog_id", blog.ID.String()),
	)

	return nil
}
