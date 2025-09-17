package service

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o servicefakes/fake_user_service.go . UserService

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (GetUserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, req http_server.PaginationRequest) ([]ListUsersResponse, int, error)
}
