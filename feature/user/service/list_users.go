package service

import (
	"context"
	"log/slog"
)

func (s *userService) ListUsers(ctx context.Context, req ListUsersRequest) ([]ListUsersResponse, int64, error) {
	s.log.Info("Listing users",
		slog.Int("page", req.Page),
		slog.Int("page_size", req.PageSize),
	)

	offset := (req.Page - 1) * req.PageSize

	users, err := s.userRepo.List(ctx, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var responses []ListUsersResponse
	for _, user := range users {
		response := ToListUsersResponse(user)
		responses = append(responses, response)
	}

	s.log.Info("Users listed successfully",
		slog.Int("count", len(responses)),
		slog.Int64("total", total),
	)

	return responses, total, nil
}
