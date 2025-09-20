package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM blogs`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		r.log.Error("Failed to count blogs",
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("%w: %w", ErrFailedToCountBlogs, err)
	}

	return count, nil
}
