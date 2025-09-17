package repository_test

import (
	"context"
	"database/sql"
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

func TestGetByIDUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogID := uuid.New()
	authorID := uuid.New()
	publishedAt := time.Now()
	expectedBlog := repository.Blog{
		ID:          blogID,
		Title:       "Test Blog",
		Content:     "Test content",
		AuthorID:    authorID,
		Status:      repository.StatusPublished,
		PublishedAt: &publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "published_at", "created_at", "updated_at"}).
		AddRow(expectedBlog.ID, expectedBlog.Title, expectedBlog.Content, expectedBlog.AuthorID, expectedBlog.Status, expectedBlog.PublishedAt, expectedBlog.CreatedAt, expectedBlog.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM blogs WHERE id = ?").
		WithArgs(blogID).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, blogID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBlog.ID, result.ID)
	assert.Equal(t, expectedBlog.Title, result.Title)
	assert.Equal(t, expectedBlog.Content, result.Content)
	assert.Equal(t, expectedBlog.AuthorID, result.AuthorID)
	assert.Equal(t, expectedBlog.Status, result.Status)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogID := uuid.New()

	// Mock the SELECT query to return no rows
	mock.ExpectQuery("SELECT (.+) FROM blogs WHERE id = ?").
		WithArgs(blogID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(ctx, blogID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
