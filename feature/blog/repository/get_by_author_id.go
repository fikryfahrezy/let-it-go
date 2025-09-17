package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (r *blogRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]Blog, error) {
	query := `
		SELECT id, title, content, author_id, status, published_at, created_at, updated_at
		FROM blogs
		WHERE author_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		slog.Error("Failed to get blogs by author ID",
			slog.String("error", err.Error()),
			slog.String("author_id", authorID.String()),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetBlogsByAuthor, err)
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		blog := Blog{}
		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.AuthorID,
			&blog.Status,
			&blog.PublishedAt,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan blog row",
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("%w: %w", ErrFailedToScanBlogRow, err)
		}
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating blog rows",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToIterateRows, err)
	}

	return blogs, nil
}
