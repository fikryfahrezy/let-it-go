package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.log.Error("Failed to list users",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToListUsers, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.log.Error("Failed to close list users rows", slog.String("error", err.Error()))
		}
	}()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			r.log.Error("Failed to scan user row",
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("%w: %w", ErrFailedToScanUserRow, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("Error iterating user rows",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToIterateRows, err)
	}

	return users, nil
}
