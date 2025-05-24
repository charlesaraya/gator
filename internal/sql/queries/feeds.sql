-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetUserFeeds :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, u.name AS user_name
FROM feeds AS f
JOIN users AS u ON f.user_id = u.id;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE url=$1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;