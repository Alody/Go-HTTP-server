// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
  token,
  user_id,
  created_at,
  updated_at,
  expires_at,
  revoked_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
`

type CreateRefreshTokenParams struct {
	Token     string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.ExecContext(ctx, createRefreshToken,
		arg.Token,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.ExpiresAt,
		arg.RevokedAt,
	)
	return err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT users.id AS user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at
FROM refresh_tokens
JOIN users ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
`

type GetUserFromRefreshTokenRow struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (GetUserFromRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i GetUserFromRefreshTokenRow
	err := row.Scan(&i.UserID, &i.ExpiresAt, &i.RevokedAt)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE token = $1
`

type RevokeRefreshTokenParams struct {
	Token     string
	RevokedAt sql.NullTime
}

func (q *Queries) RevokeRefreshToken(ctx context.Context, arg RevokeRefreshTokenParams) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, arg.Token, arg.RevokedAt)
	return err
}
