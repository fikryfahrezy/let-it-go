package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteBlogUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogID := uuid.New()

	// Mock the DELETE query
	mock.ExpectExec("DELETE FROM blogs WHERE id = ?").
		WithArgs(blogID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, blogID)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBlogNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	blogID := uuid.New()

	// Mock the DELETE query to return 0 affected rows
	mock.ExpectExec("DELETE FROM blogs WHERE id = ?").
		WithArgs(blogID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(ctx, blogID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
