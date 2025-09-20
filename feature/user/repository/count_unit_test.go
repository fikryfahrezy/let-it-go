package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	expectedCount := 5

	// Mock the COUNT query
	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(rows)

	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCountZeroUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	// Mock the COUNT query returning 0
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(rows)

	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
