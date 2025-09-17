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
		r.log.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID.String()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToUpdateUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToGetRowsAffected, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	r.log.Info("User updated successfully",
		slog.String("user_id", user.ID.String()),
	)

	return nil
}
