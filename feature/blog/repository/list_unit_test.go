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

func TestListBlogUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    repository.StatusPublished,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Blog 2",
			Content:   "Content 2",
			AuthorID:  uuid.New(),
			Status:    repository.StatusDraft,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "published_at", "created_at", "updated_at"})
	for _, blog := range blogs {
		rows.AddRow(blog.ID, blog.Title, blog.Content, blog.AuthorID, blog.Status, blog.PublishedAt, blog.CreatedAt, blog.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM blogs ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(2, 0).
		WillReturnRows(rows)

	result, err := repo.List(ctx, 2, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, blogs[0].ID, result[0].ID)
	assert.Equal(t, blogs[0].Title, result[0].Title)
	assert.Equal(t, blogs[1].ID, result[1].ID)
	assert.Equal(t, blogs[1].Title, result[1].Title)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListBlogEmptyUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	// Mock the SELECT query returning empty result
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "published_at", "created_at", "updated_at"})
	mock.ExpectQuery("SELECT (.+) FROM blogs ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(10, 0).
		WillReturnRows(rows)

	result, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
