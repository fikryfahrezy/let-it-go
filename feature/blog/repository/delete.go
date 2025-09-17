package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (r *blogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM blogs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToDeleteBlog, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToGetRowsAffected, err)
	}

	if rowsAffected == 0 {
		return ErrBlogNotFound
	}

	slog.Info("Blog deleted successfully",
		slog.String("blog_id", id.String()),
	)

	return nil
}