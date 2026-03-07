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

-- name: GetTweetCount :one
SELECT COUNT(*) FROM tweets;

-- name: GetUser :one
SELECT * 
FROM users
WHERE id = $1;

-- name: GetTweetCountByUserID :one
SELECT COUNT(*)
FROM tweets
WHERE user_id = $1;

-- name: GetRetweetCountByUserID :one
SELECT COUNT(*)
FROM retweets
WHERE user_id = $1;

-- いいね機能
-- name: CreateLike :exec
INSERT INTO likes (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateRetweet :exec
INSERT INTO retweets (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateBookmark :exec
INSERT INTO bookmarks (
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

-- name: DeleteRetweet :exec
DELETE
FROM retweets
WHERE user_id = $1 AND tweet_id = $2;

-- name: DeleteBookmark :exec
DELETE
FROM bookmarks
WHERE user_id = $1 AND tweet_id = $2;

-- GetTweetWithLikesの単体SQL
-- name: GetLikeExists :one
SELECT EXISTS (
  SELECT 1 
  FROM likes 
  WHERE user_id = $1 AND tweet_id = $2
);

-- name: GetRetweetExists :one
SELECT EXISTS (
  SELECT 1
  FROM retweets
  WHERE user_id = $1 AND tweet_id = $2
);

-- name: GetBookmarkExists :one
SELECT EXISTS (
  SELECT 1
  FROM bookmarks
  WHERE user_id = $1 AND tweet_id = $2
);

-- ツイート詳細、いいね、リツイート、ブックマーク付き
-- name: GetTweet :one
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  COUNT(DISTINCT l.id) AS like_count,
  MAX(CASE WHEN l.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_liked,
  COUNT(DISTINCT r.id) AS retweet_count,
  MAX(CASE WHEN r.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_retweeted,
  MAX(CASE WHEN b.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_bookmarked
FROM tweets t
LEFT JOIN likes l ON l.tweet_id = t.id
LEFT JOIN retweets r ON r.tweet_id = t.id
LEFT JOIN bookmarks b ON b.tweet_id = t.id
WHERE t.id = $2
GROUP BY t.id;


-- ツイート一覧、いいね、リツイート、ブックマーク付き
-- name: GetTweets :many
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  COUNT(DISTINCT l.id) AS like_count,
  MAX(CASE WHEN l.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_liked,
  COUNT(DISTINCT r.id) AS retweet_count,
  MAX(CASE WHEN r.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_retweeted,
  MAX(CASE WHEN b.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_bookmarked
FROM tweets t
LEFT JOIN likes l ON l.tweet_id = t.id
LEFT JOIN retweets r ON r.tweet_id = t.id
LEFT JOIN bookmarks b ON b.tweet_id = t.id
GROUP BY t.id
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- ユーザー詳細でのツイート一覧、いいね、リツイート、ブックマーク付き
-- name: GetTweetsByUserID :many
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  COUNT(DISTINCT l.id) AS like_count,
  MAX(CASE WHEN l.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_liked,
  COUNT(DISTINCT r.id) AS retweet_count,
  MAX(CASE WHEN r.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_retweeted,
  MAX(CASE WHEN b.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_bookmarked
FROM tweets t
LEFT JOIN likes l ON l.tweet_id = t.id
LEFT JOIN  retweets r ON r.tweet_id = t.id
LEFT JOIN bookmarks b ON b.tweet_id = t.id
WHERE t.user_id = @target_user_id::int
GROUP BY t.id
ORDER BY t.created_at DESC
LIMIT @limit_val::int OFFSET @offset_val::int;

-- 選択したユーザーがリツイートしているツイート一覧
-- name: GetRetweetedTweetsByUserID :many
SELECT
  t.id,
  t.user_id,
  t.content,
  user_retweets.created_at,
  COUNT(DISTINCT l.id) AS like_count,
  MAX(CASE WHEN l.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_liked,
  COUNT(DISTINCT all_retweets.id) AS retweet_count,
  MAX(CASE WHEN all_retweets.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_retweeted,
  MAX(CASE WHEN b.user_id = @logged_user_id::int THEN 1 ELSE 0 END)::boolean AS is_bookmarked
FROM retweets user_retweets
JOIN tweets t ON user_retweets.tweet_id = t.id
LEFT JOIN likes l ON l.tweet_id = t.id
LEFT JOIN  retweets all_retweets ON all_retweets.tweet_id = t.id
LEFT JOIN bookmarks b ON b.tweet_id = t.id
WHERE user_retweets.user_id = @target_user_id
GROUP BY t.id, user_retweets.created_at
ORDER BY user_retweets.created_at DESC
LIMIT @limit_val::int OFFSET @offset_val::int;

-- ログインしているユーザーのブックマークしたツイート一覧
-- name: GetBookmarkedTweetsByUserID :many
SELECT
  t.id,
  t.user_id,
  t.content,
  t.created_at,
  COUNT(DISTINCT l.id) AS like_count,
  MAX(CASE WHEN l.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_liked,
  COUNT(DISTINCT r.id) AS retweet_count,
  MAX(CASE WHEN r.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_retweeted,
  MAX(CASE WHEN b.user_id = $1 THEN 1 ELSE 0 END)::boolean AS is_bookmarked
FROM tweets t
LEFT JOIN likes l ON l.tweet_id = t.id
LEFT JOIN retweets r ON r.tweet_id = t.id
LEFT JOIN bookmarks b ON b.tweet_id = t.id
WHERE b.user_id = $1
GROUP BY t.id
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;