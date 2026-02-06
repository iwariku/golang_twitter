-- name: CreateUser :one
INSERT INTO users (
  mail,
  password
) VALUES (
  $1, $2
)
RETURNING *;