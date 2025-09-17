package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
)

type ListUsersResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToListUsersResponse(u repository.User) ListUsersResponse {
	return ListUsersResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
