package repository

import "context"

type BlogRepository interface {
	Create(ctx context.Context, blog Blog) error
	GetByID(ctx context.Context, id int) (Blog, error)
	GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]Blog, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]Blog, error)
	Count(ctx context.Context) (int, error)
	CountByStatus(ctx context.Context, status string) (int, error)
}