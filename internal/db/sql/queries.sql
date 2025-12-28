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

-- name: GetUserByUsername :one
SELECT * FROM flexdesk.users
WHERE username = $1 LIMIT 1;
