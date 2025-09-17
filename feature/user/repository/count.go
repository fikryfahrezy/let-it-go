package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		slog.Error("Failed to count users",
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("%w: %w", ErrFailedToCountUsers, err)
	}

	return count, nil
}