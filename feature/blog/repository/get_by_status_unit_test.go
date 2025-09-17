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

func TestGetByStatusUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	publishedBlogs := []repository.Blog{
		{
			ID:        uuid.New(),
			Title:     "Published Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    repository.StatusPublished,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "status", "published_at", "created_at", "updated_at"})
	for _, blog := range publishedBlogs {
		rows.AddRow(blog.ID, blog.Title, blog.Content, blog.AuthorID, blog.Status, blog.PublishedAt, blog.CreatedAt, blog.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM blogs WHERE status = (.+) ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(repository.StatusPublished, 10, 0).
		WillReturnRows(rows)

	result, err := repo.GetByStatus(ctx, repository.StatusPublished, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, publishedBlogs[0].ID, result[0].ID)
	assert.Equal(t, publishedBlogs[0].Status, result[0].Status)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
