package service

import "context"

type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
	GetUserByID(ctx context.Context, id int) (GetUserResponse, error)
	UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, page, pageSize int) ([]ListUsersResponse, PaginationResponse, error)
}

type PaginationResponse struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}