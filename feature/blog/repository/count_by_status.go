package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM blogs WHERE status = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		r.log.Error("Failed to count blogs by status",
			slog.String("error", err.Error()),
			slog.String("status", status),
		)
		return 0, fmt.Errorf("%w: %w", ErrFailedToCountBlogsByStatus, err)
	}

	return count, nil
}
