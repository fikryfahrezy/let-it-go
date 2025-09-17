package service

import (
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlogService_DeleteBlog_Success(t *testing.T) {
	suite := SetupBlogServiceTest()

	blogID := uuid.New()
	suite.mockRepo.DeleteReturns(nil)

	err := suite.blogService.DeleteBlog(suite.ctx, blogID)

	assert.NoError(t, err)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.DeleteCallCount())
	_, actualID := suite.mockRepo.DeleteArgsForCall(0)
	assert.Equal(t, blogID, actualID)
}

func TestBlogService_DeleteBlog_NotFound(t *testing.T) {
	suite := SetupBlogServiceTest()

	blogID := uuid.New()
	suite.mockRepo.DeleteReturns(repository.ErrBlogNotFound)

	err := suite.blogService.DeleteBlog(suite.ctx, blogID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)

	// Verify repository calls
	assert.Equal(t, 1, suite.mockRepo.DeleteCallCount())
}

