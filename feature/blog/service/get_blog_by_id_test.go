package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository/repositoryfakes"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlogService_GetBlogByID_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	blogID := uuid.New()
	authorID := uuid.New()
	expectedBlog := repository.Blog{
		ID:        blogID,
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		AuthorID:  authorID,
		Status:    "published",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.GetByIDReturns(expectedBlog, nil)

	result, err := blogService.GetBlogByID(ctx, blogID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBlog.ID, result.ID)
	assert.Equal(t, expectedBlog.Title, result.Title)
	assert.Equal(t, expectedBlog.Content, result.Content)
	assert.Equal(t, expectedBlog.AuthorID, result.AuthorID)
	assert.Equal(t, expectedBlog.Status, result.Status)
	assert.Equal(t, expectedBlog.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedBlog.UpdatedAt, result.UpdatedAt)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
	_, actualID := mockRepo.GetByIDArgsForCall(0)
	assert.Equal(t, blogID, actualID)
}

func TestBlogService_GetBlogByID_NotFound(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	blogID := uuid.New()
	mockRepo.GetByIDReturns(repository.Blog{}, repository.ErrBlogNotFound)

	result, err := blogService.GetBlogByID(ctx, blogID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)
	assert.Equal(t, service.GetBlogResponse{}, result)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.GetByIDCallCount())
}
