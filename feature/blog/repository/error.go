package repository

import "github.com/fikryfahrezy/let-it-go/pkg/app_error"

// Repository errors
var (
	// Blog not found errors
	ErrBlogNotFound = app_error.New("BLOG-BLOG_NOT_FOUND", "blog not found")

	// Database operation errors
	ErrFailedToCreateBlog       = app_error.New("BLOG-FAILED_TO_CREATE_BLOG", "failed to create blog")
	ErrFailedToGetBlog          = app_error.New("BLOG-FAILED_TO_GET_BLOG", "failed to get blog by ID")
	ErrFailedToGetBlogsByAuthor = app_error.New("BLOG-FAILED_TO_GET_BLOGS_BY_AUTHOR", "failed to get blogs by author ID")
	ErrFailedToGetBlogsByStatus = app_error.New("BLOG-FAILED_TO_GET_BLOGS_BY_STATUS", "failed to get blogs by status")
	ErrFailedToUpdateBlog       = app_error.New("BLOG-FAILED_TO_UPDATE_BLOG", "failed to update blog")
	ErrFailedToDeleteBlog       = app_error.New("BLOG-FAILED_TO_DELETE_BLOG", "failed to delete blog")
	ErrFailedToListBlogs        = app_error.New("BLOG-FAILED_TO_LIST_BLOGS", "failed to list blogs")
	ErrFailedToCountBlogs       = app_error.New("BLOG-FAILED_TO_COUNT_BLOGS", "failed to count blogs")
	ErrFailedToCountBlogsByStatus = app_error.New("BLOG-FAILED_TO_COUNT_BLOGS_BY_STATUS", "failed to count blogs by status")

	// Row scanning errors
	ErrFailedToScanBlogRow = app_error.New("BLOG-FAILED_TO_SCAN_BLOG_ROW", "failed to scan blog row")

	// Database result errors
	ErrFailedToGetLastInsertID = app_error.New("BLOG-FAILED_TO_GET_LAST_INSERT_ID", "failed to get last insert id")
	ErrFailedToGetRowsAffected = app_error.New("BLOG-FAILED_TO_GET_ROWS_AFFECTED", "failed to get rows affected")
	ErrFailedToIterateRows     = app_error.New("BLOG-FAILED_TO_ITERATE_ROWS", "error iterating blog rows")
)