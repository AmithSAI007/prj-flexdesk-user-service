-- name: CreateUser :one
INSERT INTO flexdesk.users (
  username, email, password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM flexdesk.users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM flexdesk.users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM flexdesk.users
WHERE username = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO flexdesk.refresh_tokens (
  user_id, token_hash, expires_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM flexdesk.refresh_tokens
WHERE token_hash = $1;
