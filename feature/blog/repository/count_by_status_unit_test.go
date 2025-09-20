package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountByStatusBlogUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewBlogRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	expectedCount := 3

	// Mock the COUNT query with status filter
	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM blogs WHERE status = ?").
		WithArgs(repository.StatusPublished).
		WillReturnRows(rows)

	count, err := repo.CountByStatus(ctx, repository.StatusPublished)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
