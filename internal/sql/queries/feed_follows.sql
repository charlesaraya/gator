-- name: CreateFeedFollow :one
WITH inserted AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT 
  inserted.id, 
  inserted.created_at, 
  inserted.updated_at, 
  f.name AS feed_name, 
  u.name AS user_name
FROM inserted
JOIN users AS u ON inserted.user_id = u.id
JOIN feeds AS f ON inserted.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, f.name AS feed_name
FROM feed_follows AS ff
JOIN feeds AS f ON ff.feed_id = f.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :one
WITH feed AS (
  SELECT id AS feed_id, user_id
  FROM feeds AS f 
  WHERE f.user_id = $1 AND f.url = $2
)
DELETE FROM feed_follows AS ff USING feed
WHERE ff.feed_id = feed.feed_id AND ff.user_id = feed.user_id
RETURNING ff.feed_id, ff.user_id;