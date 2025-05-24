-- name: CreatePost :one
INSERT INTO posts (id, feed_id, created_at, updated_at, title, url, description, published_at)
VALUES (
    gen_random_uuid (),
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetPostsFromUser :many
WITH userposts AS (
    SELECT ff.feed_id FROM feed_follows as ff WHERE ff.user_id = $1
)
SELECT *
FROM posts AS p
INNER JOIN userposts ON p.feed_id = userposts.feed_id
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;