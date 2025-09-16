package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM blogs`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		slog.Error("Failed to count blogs",
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("failed to count blogs: %w", err)
	}

	return count, nil
}

func (r *blogRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	query := `SELECT COUNT(*) FROM blogs WHERE status = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		slog.Error("Failed to count blogs by status",
			slog.String("error", err.Error()),
			slog.String("status", status),
		)
		return 0, fmt.Errorf("failed to count blogs by status: %w", err)
	}

	return count, nil
}