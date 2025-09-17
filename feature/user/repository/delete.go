package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToDeleteUser, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%w: %w", ErrFailedToGetRowsAffected, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	slog.Info("User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}