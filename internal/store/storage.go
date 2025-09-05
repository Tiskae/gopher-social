package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (Post, error)
		Delete(ctx context.Context, id int64) error
		UpdateOne(ctx context.Context, id int64, post *Post) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
