package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (r *userRepository) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Generate UUIDv7 for the user ID
	user.ID = uuid.Must(uuid.NewV7())

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, now, now)
	if err != nil {
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
			slog.String("email", user.Email),
		)
		return fmt.Errorf("%w: %w", ErrFailedToCreateUser, err)
	}

	// No need to get last insert ID since we're using UUIDs

	slog.Info("User created successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("email", user.Email),
	)

	return nil
}
