-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, url, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: FetchUserPosts :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.description, posts.url, posts.published_at, posts.feed_id
FROM
(SELECT * FROM feed_follows WHERE user_id = $1) AS feed_follows JOIN posts
ON feed_follows.feed_id = posts.feed_id
ORDER BY posts.created_at DESC
OFFSET $2 LIMIT $3;