-- name: GetUserStatistics :one
SELECT * FROM user_statistics
WHERE user_id = $1 LIMIT 1;

-- name: CreateUserStatistics :one
INSERT INTO user_statistics (
  user_id, total_words_learned, accuracy, total_time
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUserStatistics :one
UPDATE user_statistics
SET total_words_learned = $1, accuracy = $2, total_time = $3, updated_at = NOW()
WHERE user_id = $4
RETURNING *;

-- name: DeleteUserStatistics :exec
DELETE FROM user_statistics
WHERE user_id = $1;
