package service

import (
	"time"
	
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
)

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"min=2,max=100"`
	Email string `json:"email" validate:"email"`
}

type UpdateUserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUpdateUserResponse(u *repository.User) UpdateUserResponse {
	return UpdateUserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}