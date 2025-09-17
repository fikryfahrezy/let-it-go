package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByIDUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := repository.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Password, result.Password)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	userID := uuid.New()

	// Mock the SELECT query to return no rows
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
