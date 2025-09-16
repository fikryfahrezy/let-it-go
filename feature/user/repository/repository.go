package repository

import (
	"github.com/fikryfahrezy/let-it-go/pkg/database"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}
