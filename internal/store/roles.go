package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func (r *RolesStore) GetByName(ctx context.Context, user *User, roleName string) (bool, error) {
	query := `
		SELECT id, name, level, description FROM roles
		WHERE name = $1	
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := Role{}

	err := r.db.QueryRowContext(ctx, query, roleName).Scan(&role.ID, &role.Name, &role.Level, &role.Description)

	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
