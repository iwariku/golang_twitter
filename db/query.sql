-- name: CreateUser :one
INSERT INTO users (
  email,
  password,
  is_active,
  activation_token
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ActivateUser :exec
UPDATE users
SET is_active = true, activation_token = NULL, activated_at = CURRENT_TIMESTAMP
WHERE activation_token = $1;

-- name: GetUserByEmail :one
SELECT id, email, password, is_active
FROM users
WHERE email = $1;

-- name: CreateTweet :one
INSERT INTO tweets (
  user_id,
  content
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTweets :many
SELECT * FROM tweets
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: GetTweetCount :one
SELECT COUNT(*) FROM tweets;

-- name: GetTweet :one
SELECT id, user_id, content
FROM tweets
WHERE id = $1;

-- name: GetUser :one
SELECT * 
FROM users
WHERE id = $1;

-- name: GetTweetsByUserID :many
SELECT *
FROM tweets
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;