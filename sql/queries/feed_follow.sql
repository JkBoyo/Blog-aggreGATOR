-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
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
	SELECT 
	inserted_feed_follows.*,
	feeds.name AS feed_name,
	users.name AS user_name
FROM inserted_feed_follows
	INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
	INNER JOIN users ON inserted_feed_follows.user_id = users.id;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
	feed_follows.*,
	feeds.name AS feed_name,
	users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: RemoveFeedFollow :exec
DELETE FROM feed_follows
USING users, feeds
WHERE users.id = feed_follows.user_id 
	AND feeds.id = feed_follows.feed_id 
	AND users.name = $1 
	AND feeds.url = $2;
