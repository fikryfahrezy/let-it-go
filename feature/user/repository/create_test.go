package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	setupTest(t)

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Create(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	// Verify user was created by email since Create generates new ID
	result, err := testRepository.GetByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID == uuid.Nil {
		t.Error("Expected non-nil ID")
	}
	if result.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, result.Name)
	}
	if result.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, result.Email)
	}
	if result.Password != user.Password {
		t.Errorf("Expected password %s, got %s", user.Password, result.Password)
	}
	if result.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("Expected non-zero UpdatedAt")
	}
}

func TestCreateDuplicateEmail(t *testing.T) {
	setupTest(t)

	user1 := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword1",
	}

	user2 := repository.User{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Same email
		Password: "hashedpassword2",
	}

	err := testRepository.Create(context.Background(), user1)
	if err != nil {
		t.Fatal(err)
	}

	err = testRepository.Create(context.Background(), user2)
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}