package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yutaronakayama/otelsql-trace-test/domain"
)

func selectUsers(ctx context.Context, db *sqlx.DB) ([]domain.User, error) {
	var users []domain.User
	err := db.SelectContext(ctx, &users, "SELECT user_id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}
