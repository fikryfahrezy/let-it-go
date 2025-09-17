package service

import (
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_ListUsers_Success(t *testing.T) {
	suite := SetupUserServiceTest()

	expectedUsers := []repository.User{
		{
			ID:        uuid.New(),
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "hashedpassword1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "hashedpassword2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	suite.mockRepo.ListReturns(expectedUsers, nil)
	suite.mockRepo.CountReturns(2, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	result, totalCount, err := suite.userService.ListUsers(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, totalCount)

	// Verify first user
	assert.Equal(t, expectedUsers[0].ID, result[0].ID)
	assert.Equal(t, expectedUsers[0].Name, result[0].Name)
	assert.Equal(t, expectedUsers[0].Email, result[0].Email)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	_, limit, offset := suite.mockRepo.ListArgsForCall(0)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset) // (page-1) * pageSize = (1-1) * 10 = 0

	assert.Equal(t, 1, suite.mockRepo.CountCallCount())
}

func TestUserService_ListUsers_WithCustomPagination(t *testing.T) {
	suite := SetupUserServiceTest()

	expectedUsers := []repository.User{
		{
			ID:        uuid.New(),
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "hashedpassword1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	suite.mockRepo.ListReturns(expectedUsers, nil)
	suite.mockRepo.CountReturns(25, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     3,
		PageSize: 5,
	}

	result, totalCount, err := suite.userService.ListUsers(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 25, totalCount)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	_, limit, offset := suite.mockRepo.ListArgsForCall(0)
	assert.Equal(t, 5, limit)
	assert.Equal(t, 10, offset) // (page-1) * pageSize = (3-1) * 5 = 10
}

func TestUserService_ListUsers_EmptyResult(t *testing.T) {
	suite := SetupUserServiceTest()

	suite.mockRepo.ListReturns([]repository.User{}, nil)
	suite.mockRepo.CountReturns(0, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	result, totalCount, err := suite.userService.ListUsers(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, totalCount)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	assert.Equal(t, 1, suite.mockRepo.CountCallCount())
}

