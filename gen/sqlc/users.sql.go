// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    user_id,
    email,
    password
) VALUES ($1, $2, $3)
`

type CreateUserParams struct {
	UserID   uuid.UUID
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser, arg.UserID, arg.Email, arg.Password)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT email, password FROM users 
    WHERE email = $1 LIMIT 1
`

type GetUserByEmailRow struct {
	Email    string
	Password string
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.Email, &i.Password)
	return i, err
}
