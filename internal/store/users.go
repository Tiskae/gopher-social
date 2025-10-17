package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email,omitempty"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at,omitempty"`
	IsActive  bool     `json:"is_active,omitempty"`
	RoleID    int      `json:"role_id"`
	Role      Role     `json:"role"`
}

type Password struct {
	text string
	hash []byte
}

func (p *Password) Set(value string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = value
	p.hash = hash

	return nil
}

func (p *Password) CompareHash(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))

	return err
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4)) RETURNING id, created_at, role_id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := s.db.
		QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash, role).
		Scan(&user.ID, &user.CreatedAt, &user.RoleID)

	user.Role.ID = user.RoleID

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.created_at, r.id, r.name, r.level FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND is_active = true
	`

	user := User{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return user, ErrNotFound
		default:
			return user, err
		}
	}

	return user, nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.created_at, r.id, r.name, r.level FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.username = $1 AND is_active = true
	`

	user := User{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return user, ErrNotFound
		default:
			return user, err
		}
	}

	return user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invite
		err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find the user that the token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		// update the active status of the user
		err = s.activateUser(ctx, tx, user)
		if err != nil {
			return err
		}

		// cleanup the invitations
		err = s.cleanupInvitations(ctx, tx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.created_at, u.is_active FROM users u
		JOIN user_invitations ui
			ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.CreatedAt,
			&user.IsActive)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) activateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET is_active = TRUE
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	user.IsActive = true

	return nil
}

func (s *UserStore) cleanupInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return nil
	}

	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.cleanupInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}
