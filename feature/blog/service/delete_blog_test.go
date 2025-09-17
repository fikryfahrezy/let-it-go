package service_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository/repositoryfakes"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlogService_DeleteBlog_Success(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	blogID := uuid.New()
	mockRepo.DeleteReturns(nil)

	err := blogService.DeleteBlog(ctx, blogID)

	assert.NoError(t, err)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.DeleteCallCount())
	_, actualID := mockRepo.DeleteArgsForCall(0)
	assert.Equal(t, blogID, actualID)
}

func TestBlogService_DeleteBlog_NotFound(t *testing.T) {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := service.NewBlogService(logger.NewDiscardLogger(), mockRepo)
	ctx := context.Background()

	blogID := uuid.New()
	mockRepo.DeleteReturns(repository.ErrBlogNotFound)

	err := blogService.DeleteBlog(ctx, blogID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)

	// Verify repository calls
	assert.Equal(t, 1, mockRepo.DeleteCallCount())
}
