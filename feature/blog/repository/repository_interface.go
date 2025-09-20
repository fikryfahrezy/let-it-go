package repository

//counterfeiter:generate -o repositoryfakes/fake_blog_repository.go . BlogRepository

import (
	"context"

	"github.com/google/uuid"
)

type BlogRepository interface {
	Create(ctx context.Context, blog Blog) error
	GetByID(ctx context.Context, id uuid.UUID) (Blog, error)
	GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]Blog, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]Blog, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}
