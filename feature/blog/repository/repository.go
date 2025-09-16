package repository

import (
	"github.com/fikryfahrezy/let-it-go/pkg/database"
)

type blogRepository struct {
	db *database.DB
}

func NewBlogRepository(db *database.DB) *blogRepository {
	return &blogRepository{
		db: db,
	}
}
