package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		r.log.Error("Failed to count users",
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("%w: %w", ErrFailedToCountUsers, err)
	}

	return count, nil
}
