package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func (r *userRepository) GetByID(ctx context.Context, id int) (User, error) {
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
			return User{}, fmt.Errorf("user not found")
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}