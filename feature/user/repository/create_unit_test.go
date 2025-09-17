package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	// Mock the INSERT query
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), user.Name, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateDuplicateEmailUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	// Mock the INSERT query to return a duplicate entry error
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), user.Name, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone) // Simulate a database error

	err = repo.Create(ctx, user)
	assert.Error(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
