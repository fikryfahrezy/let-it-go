package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.log.Info("Deleting user",
		slog.String("user_id", id.String()),
	)

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.log.Info("User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}
