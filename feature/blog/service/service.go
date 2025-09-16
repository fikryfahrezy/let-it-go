package service

import "github.com/fikryfahrezy/let-it-go/feature/blog/repository"

type blogService struct {
	blogRepo repository.BlogRepository
}

func NewBlogService(blogRepo repository.BlogRepository) BlogService {
	return &blogService{
		blogRepo: blogRepo,
	}
}