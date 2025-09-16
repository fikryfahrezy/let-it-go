package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM blogs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete blog",
			slog.String("error", err.Error()),
			slog.Int("blog_id", id),
		)
		return fmt.Errorf("failed to delete blog: %w", err)
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

	slog.Info("Blog deleted successfully",
		slog.Int("blog_id", id),
	)

	return nil
}