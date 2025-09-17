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

func TestUserService_UpdateUser_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	existingUser := repository.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	mockRepo.GetByIDReturns(existingUser, nil)
	mockRepo.UpdateReturns(nil)

	req := service.UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	result, err := userService.UpdateUser(ctx, userID, req)

	assert.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, existingUser.CreatedAt, result.CreatedAt)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
	assert.Equal(t, 1, mockRepo.UpdateCallCount())
	_, actualUser := mockRepo.UpdateArgsForCall(0)
	assert.Equal(t, req.Name, actualUser.Name)
	assert.Equal(t, req.Email, actualUser.Email)
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.GetByIDReturns(repository.User{}, repository.ErrUserNotFound)

	req := service.UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	result, err := userService.UpdateUser(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
	assert.Equal(t, service.UpdateUserResponse{}, result)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
	assert.Equal(t, 0, mockRepo.UpdateCallCount()) // Update should not be called
}
