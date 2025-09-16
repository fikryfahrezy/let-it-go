package service

import "github.com/fikryfahrezy/let-it-go/feature/user/repository"

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}


