
-- name: AddFeed :one 
INSERT INTO feeds (id, created_at, updated_at, feed_name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT f.feed_name, f.url, u.name 
FROM feeds as f
LEFT JOIN users as u
ON f.user_id = u.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
  )
RETURNING *
)
SELECT inserted_feed_follow.*,
feeds.feed_name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE feeds.url = $1;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.feed_name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
