package service

import (
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_DeleteUser_Success(t *testing.T) {
	suite := SetupUserServiceTest()

	userID := uuid.New()
	suite.mockRepo.DeleteReturns(nil)

	err := suite.userService.DeleteUser(suite.ctx, userID)

	assert.NoError(t, err)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.DeleteCallCount())
	_, actualID := suite.mockRepo.DeleteArgsForCall(0)
	assert.Equal(t, userID, actualID)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	suite := SetupUserServiceTest()

	userID := uuid.New()
	suite.mockRepo.DeleteReturns(repository.ErrUserNotFound)

	err := suite.userService.DeleteUser(suite.ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.DeleteCallCount())
}

