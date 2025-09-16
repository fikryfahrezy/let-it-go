package service

import "github.com/fikryfahrezy/let-it-go/feature/user/repository"

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}
