package domain

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type User struct {
	UserID int    `db:"user_id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
}

type UserRepository interface {
	SelectUsers(ctx context.Context, db *sqlx.DB) ([]User, error)
}
