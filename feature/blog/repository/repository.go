package repository

import (
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
)

type blogRepository struct {
	db  *database.DB
	log *slog.Logger
}

func NewBlogRepository(log *slog.Logger, db *database.DB) *blogRepository {
	return &blogRepository{
		db:  db,
		log: log,
	}
}
