-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
  token,
  user_id,
  created_at,
  updated_at,
  expires_at,
  revoked_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.id AS user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at
FROM refresh_tokens
JOIN users ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;