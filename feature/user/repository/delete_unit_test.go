package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	userID := uuid.New()

	// Mock the DELETE query
	mock.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, userID)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	userID := uuid.New()

	// Mock the DELETE query to return 0 affected rows
	mock.ExpectExec("DELETE FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
