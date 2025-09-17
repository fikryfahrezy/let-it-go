package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository/repositoryfakes"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := repository.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.GetByIDReturns(expectedUser, nil)

	result, err := userService.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedUser.UpdatedAt, result.UpdatedAt)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
	_, actualID := mockRepo.GetByIDArgsForCall(0)
	assert.Equal(t, userID, actualID)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.GetByIDReturns(repository.User{}, repository.ErrUserNotFound)

	result, err := userService.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
	assert.Equal(t, service.GetUserResponse{}, result)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
}
