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
SELECT *
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

-- name: GetTweetsByUserIDWithLikes :many
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS like_count,
  EXISTS (
    SELECT 1
    FROM likes l
    WHERE l.tweet_id = t.id AND l.user_id = @viewer_user_id::int
  ) AS is_liked
FROM tweets t
WHERE t.user_id = @target_user_id::int
ORDER BY t.created_at DESC
LIMIT @limit_val::int OFFSET @offset_val::int;


-- name: GetTweetCountByUserID :one
SELECT COUNT(*)
FROM tweets
WHERE user_id = $1;

-- いいね機能
-- name: CreateLike :one
INSERT INTO likes (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteLike :exec
DELETE 
FROM likes
WHERE user_id = $1 AND tweet_id = $2;

-- GetTweetWithLikesの単体SQL
-- name: GetLikeExists :one
SELECT EXISTS (
  SELECT 1 
  FROM likes 
  WHERE user_id = $1 AND tweet_id = $2
);

-- GetTweetWithLikesの単体SQL
-- name: GetLikeCountByTweetID :one
SELECT COUNT(*)
FROM likes
WHERE tweet_id = $1;

-- name: GetTweetWithLikes :one
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS like_count,
  EXISTS (
    SELECT 1
    FROM likes l
    WHERE l.tweet_id = t.id AND l.user_id = $1
  ) AS is_liked
FROM tweets t
WHERE t.id = $2;

-- name: GetTweetsWithLikes :many
SELECT 
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS like_count,
  EXISTS (
    SELECT 1
    FROM likes l
    WHERE l.tweet_id = t.id AND l.user_id = $1
  ) AS is_liked
FROM tweets t
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;