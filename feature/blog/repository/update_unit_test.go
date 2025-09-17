package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBlogUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogID := uuid.New()
	authorID := uuid.New()
	blog := repository.Blog{
		ID:        blogID,
		Title:     "Updated Blog Title",
		Content:   "Updated content",
		AuthorID:  authorID,
		Status:    repository.StatusPublished,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock the UPDATE query - matches the actual query parameters
	mock.ExpectExec("UPDATE blogs SET").
		WithArgs(blog.Title, blog.Content, blog.Status, blog.PublishedAt, sqlmock.AnyArg(), blog.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, blog)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBlogNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blog := repository.Blog{
		ID:       uuid.New(),
		Title:    "Non-existent Blog",
		Content:  "Content",
		AuthorID: uuid.New(),
		Status:   repository.StatusDraft,
	}

	// Mock the UPDATE query to return 0 affected rows
	mock.ExpectExec("UPDATE blogs SET").
		WithArgs(blog.Title, blog.Content, blog.Status, blog.PublishedAt, sqlmock.AnyArg(), blog.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(ctx, blog)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
