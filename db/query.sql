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

-- name: CreateRetweet :one
INSERT INTO retweets (
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

-- リツイート機能
-- ツイート詳細、いいね、リツイート付き(リツイートができたら巻き替え)
-- name: GetTweetsWithLikeAndRetweet :one
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
  ) AS is_liked,
  (SELECT COUNT(*) FROM retweets r WHERE r.tweet_id = t.id ) AS retweet_count,
  EXISTS (
    SELECT 1
    FROM retweets r
    WHERE r.tweet_id = t.id AND r.user_id = $1
  ) AS is_retweeted
FROM tweets t
WHERE t.id = $2;

-- ツイート一覧、いいね、リツイート付き(リツイートができたら巻き替え)
-- name: GetTweetsWithRetweet :many
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
  ) AS is_liked,
  (SELECT COUNT(*) FROM retweets r WHERE r.tweet_id = t.id) AS retweet_count,
  EXISTS (
    SELECT 1
    FROM retweets r
    WHERE r.tweet_id = t.id AND r.user_id = $1
  ) AS is_retweeted
FROM tweets t
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- ユーザー詳細でのツイート一覧、いいね、リツイート付き(リツイートができたら巻き替え)
-- name: GetTweetsByUserIDWithRetweet :many
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
  ) AS is_liked,
  (SELECT COUNT(*) FROM retweets r WHERE r.tweet_id = t.id) AS retweet_count,
  EXISTS (
    SELECT 1
    FROM retweets r
    WHERE r.tweet_id = t.id AND r.user_id = @viewer_user_id::int
  ) AS is_retweeted
FROM tweets t
WHERE t.user_id = @target_user_id::int
ORDER BY t.created_at DESC
LIMIT @limit_val::int OFFSET @offset_val::int;

-- 選択したユーザーがリツイートしているツイート一覧
-- name: GetRetweetedTweetsByUserID :many
SELECT
  t.id,
  t.user_id,
  t.content,
  user_retweets.created_at,
  (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS like_count,
  EXISTS (
    SELECT 1
    FROM likes l
    WHERE l.tweet_id = t.id AND l.user_id = $1
  ) AS is_liked,

  (SELECT COUNT(*) FROM retweets viewer_retweet WHERE viewer_retweet.tweet_id = t.id ) AS retweet_count,
  EXISTS (
    SELECT 1
    FROM retweets viewer_retweet
    WHERE viewer_retweet.tweet_id = t.id AND viewer_retweet.user_id = $1
  ) AS is_retweeted
FROM retweets user_retweets
JOIN tweets t ON user_retweets.tweet_id = t.id
WHERE user_retweets.user_id = $2
ORDER BY user_retweets.created_at DESC
LIMIT $3 OFFSET $4;