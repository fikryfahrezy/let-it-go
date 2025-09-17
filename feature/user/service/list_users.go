package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
)

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]ListUsersResponse, PaginationResponse, error) {
	slog.Info("Listing users",
		slog.Int("page", page),
		slog.Int("page_size", pageSize),
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		slog.Error("Failed to list users",
			slog.String("error", err.Error()),
		)
		return nil, PaginationResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToListUsers, err)
	}

	total, err := s.userRepo.Count(ctx)
	if err != nil {
		slog.Error("Failed to count users",
			slog.String("error", err.Error()),
		)
		return nil, PaginationResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToCountUsers, err)
	}

	var responses []ListUsersResponse
	for _, user := range users {
		response := ToListUsersResponse(user)
		responses = append(responses, response)
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	pagination := PaginationResponse{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}

	slog.Info("Users listed successfully",
		slog.Int("count", len(responses)),
		slog.Int("total", total),
	)

	return responses, pagination, nil
}
