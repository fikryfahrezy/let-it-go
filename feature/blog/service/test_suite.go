package service

import (
	"context"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository/repositoryfakes"
)

// BlogServiceTestSuite provides common setup for blog service tests
type BlogServiceTestSuite struct {
	mockRepo    *repositoryfakes.FakeBlogRepository
	blogService BlogService
	ctx         context.Context
}

// SetupBlogServiceTest creates a new test suite instance
func SetupBlogServiceTest() *BlogServiceTestSuite {
	mockRepo := &repositoryfakes.FakeBlogRepository{}
	blogService := NewBlogService(mockRepo)
	
	return &BlogServiceTestSuite{
		mockRepo:    mockRepo,
		blogService: blogService,
		ctx:         context.Background(),
	}
}

// Helper functions for common assertions and setup can be added here