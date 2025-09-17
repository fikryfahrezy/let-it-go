package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return User{}, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	return user, nil
}