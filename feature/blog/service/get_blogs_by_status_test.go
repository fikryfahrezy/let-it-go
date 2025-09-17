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

func TestBlogService_GetBlogsByStatus_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	status := "published"
	expectedBlogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Published Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Published Blog 2",
			Content:   "Content 2",
			AuthorID:  uuid.New(),
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.GetByStatusReturns(expectedBlogs, nil)
	mockRepo.CountByStatusReturns(2, nil)

	paginationReq := http_server.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	result, totalCount, err := blogService.GetBlogsByStatus(ctx, status, paginationReq)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, totalCount)

	// Verify first blog
	assert.Equal(t, expectedBlogs[0].ID, result[0].ID)
	assert.Equal(t, expectedBlogs[0].Status, result[0].Status)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByStatusCallCount())
	_, actualStatus, limit, offset := mockRepo.GetByStatusArgsForCall(0)
	assert.Equal(t, status, actualStatus)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset)

	assert.Equal(t, 1, mockRepo.CountByStatusCallCount())
	_, actualStatusForCount := mockRepo.CountByStatusArgsForCall(0)
	assert.Equal(t, status, actualStatusForCount)
}
