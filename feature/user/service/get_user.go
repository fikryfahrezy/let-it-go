package service

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *userService) GetUserByID(ctx context.Context, id int) (GetUserResponse, error) {
	slog.Info("Getting user by ID",
		slog.Int("user_id", id),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return GetUserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		slog.Warn("User not found",
			slog.Int("user_id", id),
		)
		return GetUserResponse{}, fmt.Errorf("user not found")
	}

	response := ToGetUserResponse(user)
	return response, nil
}