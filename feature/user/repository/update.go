package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (r *userRepository) Update(ctx context.Context, user User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, now, user.ID)
	if err != nil {
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.Int("user_id", user.ID),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	slog.Info("User updated successfully",
		slog.Int("user_id", user.ID),
	)

	return nil
}
