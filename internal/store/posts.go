package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	UserID    int64    `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.
		QueryRowContext(
			ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserID).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int) (*Post, error) {
	query := `
		SELECT id, content,  title, tags, user_id, created_at, updated_at FROM posts WHERE posts.id = $1
	`

	var post *Post = &Post{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.UserID, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		return &Post{}, err
	}

	return post, nil
}
