package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)

	if err != nil {
		pqErr, ok := err.(*pq.Error)

		if ok && pqErr.Code == "23505" { // conflict error
			return ErrConflict
		} else if pqErr.Code == "23503" { // foreign key violation error
			return ErrNotFound
		} else {
			return err
		}
	}

	return nil
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row, err := s.db.ExecContext(ctx, query, userID, followerID)

	if err != nil {
		return err
	}

	rowsDeleted, err := row.RowsAffected()
	if err != nil {
		return err
	}

	// no following found, hence nothing got deleted
	if rowsDeleted == 0 {
		return ErrNotFound
	}

	return err
}
