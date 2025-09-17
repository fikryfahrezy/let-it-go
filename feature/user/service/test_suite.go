package service

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository/repositoryfakes"
)

// UserServiceTestSuite provides common setup for user service tests
type UserServiceTestSuite struct {
	mockRepo    *repositoryfakes.FakeUserRepository
	userService UserService
	ctx         context.Context
}

// SetupUserServiceTest creates a new test suite instance
func SetupUserServiceTest() *UserServiceTestSuite {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := NewUserService(mockRepo)
	
	return &UserServiceTestSuite{
		mockRepo:    mockRepo,
		userService: userService,
		ctx:         context.Background(),
	}
}

// Helper functions for common assertions and setup can be added here