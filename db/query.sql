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

-- フォロー関連
-- name: CreateFollow :exec
INSERT INTO follows (
  follower_id,
  following_id
) VALUES (
  $1, $2
);

-- name: DeleteFollow :exec
DELETE
FROM follows
WHERE follower_id = $1 AND following_id = $2;

-- name: GetFollowExists :one
SELECT EXISTS (
  SELECT 1
  FROM follows
  WHERE follower_id = $1 AND following_id = $2
);

-- name: GetFollowingCount :one
SELECT COUNT(*)
FROM follows
WHERE follower_id = $1;

-- フォロー一覧で閲覧
-- name: GetFollowings :many
SELECT
  u.id,
  u.user_name,
  u.nick_name,
  u.self_introduction,
  u.profile_image,
  EXISTS (
    SELECT 1
    FROM follows check_f
    WHERE check_f.follower_id = $1 AND check_f.following_id = u.id
  ) AS is_followed
FROM follows f
INNER JOIN users u ON f.following_id = u.id
WHERE f.follower_id = $2
ORDER BY f.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetFollowerCount :one
SELECT COUNT(*)
FROM follows
WHERE following_id = $1;

--  フォロワー一覧
-- name: GetFollowers :many
SELECT
  u.id,
  u.user_name,
  u.nick_name,
  u.self_introduction,
  u.profile_image,
  EXISTS (
    SELECT 1
    FROM follows check_f
    WHERE check_f.follower_id = $1 AND check_f.following_id = u.id
  ) AS is_followed
FROM follows f
INNER JOIN users u ON f.follower_id = u.id
WHERE f.following_id = $2
ORDER BY f.created_at DESC
LIMIT $3 OFFSET $4;

-- DM機能
-- グループ作成を実装する
-- :oneにするのはこの作成したgroupsテーブルのidを別のテーブルで使用するため。必要がない場合は:execに変更する
-- グループ作成の時にグループ名を入力するイメージ
-- 必要なデータ、誰が作ったか:user_id、グループ名:name
-- name: CreateGroup :one
INSERT INTO dm_groups (
  name
) VALUES (
  $1
)
RETURNING *;

-- グループが作成されたら、ログインしているユーザーとグループidを使って作成されたグループに自分を入れる
-- これはdm_groupsではnameカラムしか持たず、ユーザー情報はdm_group_membersに入れるという設計にしているため
-- name: AddMemberToGroup :one
INSERT INTO dm_group_members (
  user_id,
  dm_group_id
) VALUES (
  $1, $2
)
RETURNING *;

-- グループでメッセージを投稿できるようにする
-- user_idはCookieにセットしあるsessionIDを使う
-- 必要なデータ、誰が: user_id,どこに: dm_group_id、何を: message
-- name: CreateMessage :one
INSERT INTO dm_messages (
  user_id,
  dm_group_id,
  message
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- グループ内のメッセージを参照できるようにメッセージ一覧を実装する
-- 必要なデータ、誰の:user_id、メッセージか: message どこのグループに所属しているか?: dm_group_id = $1;
-- name: GetMessagesByGruopID :many
SELECT
  user_id,
  message
FROM dm_messages
WHERE dm_group_id = $1;

-- グループの一覧を参照できるようにする
-- WHEREがないと自分の所属しているグループ以外も表示されてしまう。
-- 別名のエイリアスについて質問する
-- name: GetGroups :many
SELECT
  dm_group_members.user_id,
  dm_group_members.dm_group_id,
  dm_groups.name
FROM dm_group_members
JOIN dm_groups  ON dm_group_members.dm_group_id = dm_groups.id
WHERE dm_group_members.user_id = $1;