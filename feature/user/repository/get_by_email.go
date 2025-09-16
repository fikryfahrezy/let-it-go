package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
		slog.Error("Failed to get user by email",
			slog.String("error", err.Error()),
			slog.String("email", email),
		)
		return User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}