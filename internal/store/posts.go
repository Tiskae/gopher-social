// Package store provides data access and storage functionality for posts and related entities.
package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.tags,
			u.username,
			COUNT(c.id) comments_count
		FROM
			posts p
			LEFT JOIN comments c ON p.id = c.post_id
			LEFT JOIN users u ON u.id = p.user_id
			INNER JOIN followers f ON p.user_id = f.follower_id
			OR p.user_id = $1
			WHERE f.user_id = $1 OR p.user_id = $1
		GROUP BY
			p.id,u.id
		ORDER BY p.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var feedPosts []PostWithMetadata

	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		return feedPosts, err
	}

	defer rows.Close()

	for rows.Next() {
		post := PostWithMetadata{}

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentsCount)

		if err != nil {
			return feedPosts, err
		}

		feedPosts = append(feedPosts, post)
	}

	return feedPosts, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.
		QueryRowContext(
			ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserID).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Version)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (Post, error) {
	query := `
		SELECT id, content,  title, tags, user_id, created_at, updated_at, version
		FROM posts WHERE posts.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post = Post{}

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID, &post.Content,
		&post.Title, pq.Array(&post.Tags),
		&post.UserID, &post.CreatedAt, &post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return post, ErrNotFound
		default:
			return post, err
		}
	}

	return post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE from posts WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsDeleted == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) UpdateOne(ctx context.Context, postID int64, updatedPost *Post) error {
	query := `
		UPDATE posts
		SET
			title = $2,
			content = $3,
			tags = $4,
			version = version + 1
		WHERE id = $1 AND version = $5
		RETURNING title, content, tags, user_id, created_at, updated_at, version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx, query, postID, updatedPost.Title, updatedPost.Content, pq.Array(updatedPost.Tags), updatedPost.Version).
		Scan(&updatedPost.Title, &updatedPost.Content, pq.Array(&updatedPost.Tags),
			&updatedPost.UserID, &updatedPost.CreatedAt, &updatedPost.UpdatedAt, &updatedPost.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
