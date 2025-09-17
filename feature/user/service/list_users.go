package service

import (
	"context"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
)

func (s *userService) ListUsers(ctx context.Context, req http_server.PaginationRequest) ([]ListUsersResponse, int, error) {
	slog.Info("Listing users",
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


	slog.Info("Users listed successfully",
		slog.Int("count", len(responses)),
		slog.Int("total", total),
	)

	return responses, total, nil
}
