package service

import "github.com/fikryfahrezy/let-it-go/pkg/app_error"

// Business logic errors (service-specific only)
var (
	// Blog status errors (business logic specific)
	ErrInvalidBlogStatus = app_error.New("BLOG-INVALID_BLOG_STATUS", "invalid blog status")
	ErrBlogAlreadyPublished = app_error.New("BLOG-BLOG_ALREADY_PUBLISHED", "blog is already published")
	ErrBlogAlreadyArchived = app_error.New("BLOG-BLOG_ALREADY_ARCHIVED", "blog is already archived")

	// Service-specific operation errors
	ErrFailedToPublishBlog = app_error.New("BLOG-FAILED_TO_PUBLISH_BLOG", "failed to publish blog")
	ErrFailedToArchiveBlog = app_error.New("BLOG-FAILED_TO_ARCHIVE_BLOG", "failed to archive blog")
)