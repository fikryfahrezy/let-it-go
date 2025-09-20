package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository/repositoryfakes"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlogService_GetBlogsByAuthor_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	authorID := uuid.New()
	expectedBlogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  authorID,
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Blog 2",
			Content:   "Content 2",
			AuthorID:  authorID,
			Status:    "draft",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.GetByAuthorIDReturns(expectedBlogs, nil)
	mockRepo.CountReturns(2, nil)

	paginationReq := service.GetBlogsByAuthorRequest{
		PaginationRequest: http_server.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	result, totalCount, err := blogService.GetBlogsByAuthor(ctx, authorID, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), totalCount)

	// Verify first blog
	assert.Equal(t, expectedBlogs[0].ID, result[0].ID)
	assert.Equal(t, expectedBlogs[0].Title, result[0].Title)
	assert.Equal(t, expectedBlogs[0].Content, result[0].Content)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByAuthorIDCallCount())
	_, actualAuthorID, limit, offset := mockRepo.GetByAuthorIDArgsForCall(0)
	assert.Equal(t, authorID, actualAuthorID)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset) // (page-1) * pageSize = (1-1) * 10 = 0

	assert.Equal(t, 1, mockRepo.CountCallCount())
}
