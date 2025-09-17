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

func TestGetByEmailUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	email := "john@example.com"
	expectedUser := repository.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     email,
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	result, err := repo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Password, result.Password)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmailNotFoundUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	email := "nonexistent@example.com"

	// Mock the SELECT query to return no rows
	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByEmail(ctx, email)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
