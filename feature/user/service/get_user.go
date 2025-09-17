package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserResponse, error) {
	s.log.Info("Getting user by ID",
		slog.String("user_id", id.String()),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return GetUserResponse{}, err
	}

	response := ToGetUserResponse(user)
	return response, nil
}
