package repository_test

import (
	"context"
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

func TestListUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	users := []repository.User{
		{
			ID:        uuid.New(),
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "password1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "password2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"})
	for _, user := range users {
		rows.AddRow(user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(2, 0).
		WillReturnRows(rows)

	result, err := repo.List(ctx, 2, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[0].Name, result[0].Name)
	assert.Equal(t, users[1].ID, result[1].ID)
	assert.Equal(t, users[1].Name, result[1].Name)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListEmptyUnit(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	// nolint:errcheck
	defer sqlDB.Close()

	db := &database.DB{DB: sqlDB}
	repo := repository.NewUserRepository(logger.NewDiscardLogger(), db)
	ctx := context.Background()

	// Mock the SELECT query returning empty result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"})
	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY created_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(10, 0).
		WillReturnRows(rows)

	result, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
