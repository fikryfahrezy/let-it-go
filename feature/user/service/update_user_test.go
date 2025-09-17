package service

import (
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_UpdateUser_Success(t *testing.T) {
	suite := SetupUserServiceTest()

	userID := uuid.New()
	existingUser := repository.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	suite.mockRepo.GetByIDReturns(existingUser, nil)
	suite.mockRepo.UpdateReturns(nil)

	req := UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	result, err := suite.userService.UpdateUser(suite.ctx, userID, req)

	assert.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, existingUser.CreatedAt, result.CreatedAt)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.GetByIDCallCount())
	assert.Equal(t, 1, suite.mockRepo.UpdateCallCount())
	_, actualUser := suite.mockRepo.UpdateArgsForCall(0)
	assert.Equal(t, req.Name, actualUser.Name)
	assert.Equal(t, req.Email, actualUser.Email)
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	suite := SetupUserServiceTest()

	userID := uuid.New()
	suite.mockRepo.GetByIDReturns(repository.User{}, repository.ErrUserNotFound)

	req := UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	result, err := suite.userService.UpdateUser(suite.ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
	assert.Equal(t, UpdateUserResponse{}, result)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.GetByIDCallCount())
	assert.Equal(t, 0, suite.mockRepo.UpdateCallCount()) // Update should not be called
}

