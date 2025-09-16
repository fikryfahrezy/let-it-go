package service

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	slog.Info("Deleting user",
		slog.Int("user_id", id),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.ID == 0 {
		slog.Warn("User not found",
			slog.Int("user_id", id),
		)
		return fmt.Errorf("user not found")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	slog.Info("User deleted successfully",
		slog.Int("user_id", id),
	)

	return nil
}
