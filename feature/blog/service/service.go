package service

import (
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

type blogService struct {
	blogRepo repository.BlogRepository
	log      *slog.Logger
}

func NewBlogService(log *slog.Logger, blogRepo repository.BlogRepository) *blogService {
	return &blogService{
		blogRepo: blogRepo,
		log:      log,
	}
}
