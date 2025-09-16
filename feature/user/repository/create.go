package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (r *userRepository) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, now, now)
	if err != nil {
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
			slog.String("email", user.Email),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert id",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = int(id)

	slog.Info("User created successfully",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email),
	)

	return nil
}
