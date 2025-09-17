package service_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository/repositoryfakes"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_DeleteUser_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.DeleteReturns(nil)

	err := userService.DeleteUser(ctx, userID)

	assert.NoError(t, err)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.DeleteCallCount())
	_, actualID := mockRepo.DeleteArgsForCall(0)
	assert.Equal(t, userID, actualID)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	mockRepo := &repositoryfakes.FakeUserRepository{}
	userService := service.NewUserService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	mockRepo.DeleteReturns(repository.ErrUserNotFound)

	err := userService.DeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.DeleteCallCount())
}
