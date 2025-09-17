package service

import (
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlogService_ListBlogs_Success(t *testing.T) {
	suite := SetupBlogServiceTest()

	expectedBlogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Blog 2",
			Content:   "Content 2",
			AuthorID:  uuid.New(),
			Status:    "draft",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	suite.mockRepo.ListReturns(expectedBlogs, nil)
	suite.mockRepo.CountReturns(2, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	result, totalCount, err := suite.blogService.ListBlogs(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, totalCount)

	// Verify first blog
	assert.Equal(t, expectedBlogs[0].ID, result[0].ID)
	assert.Equal(t, expectedBlogs[0].Title, result[0].Title)
	assert.Equal(t, expectedBlogs[0].Content, result[0].Content)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	_, limit, offset := suite.mockRepo.ListArgsForCall(0)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset) // (page-1) * pageSize = (1-1) * 10 = 0

	assert.Equal(t, 1, suite.mockRepo.CountCallCount())
}

func TestBlogService_ListBlogs_WithCustomPagination(t *testing.T) {
	suite := SetupBlogServiceTest()

	expectedBlogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	suite.mockRepo.ListReturns(expectedBlogs, nil)
	suite.mockRepo.CountReturns(25, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     3,
		PageSize: 5,
	}

	result, totalCount, err := suite.blogService.ListBlogs(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 25, totalCount)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	_, limit, offset := suite.mockRepo.ListArgsForCall(0)
	assert.Equal(t, 5, limit)
	assert.Equal(t, 10, offset) // (page-1) * pageSize = (3-1) * 5 = 10
}

func TestBlogService_ListBlogs_EmptyResult(t *testing.T) {
	suite := SetupBlogServiceTest()

	suite.mockRepo.ListReturns([]repository.Blog{}, nil)
	suite.mockRepo.CountReturns(0, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	result, totalCount, err := suite.blogService.ListBlogs(suite.ctx, paginationReq)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, totalCount)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.ListCallCount())
	assert.Equal(t, 1, suite.mockRepo.CountCallCount())
}

