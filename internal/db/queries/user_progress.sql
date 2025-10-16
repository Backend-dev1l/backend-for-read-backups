-- name: GetUserProgress :one
SELECT * FROM user_progress
WHERE id = $1 LIMIT 1;

-- name: GetUserProgressByUserAndWord :one
SELECT * FROM user_progress
WHERE user_id = $1 AND word_id = $2 LIMIT 1;

-- name: ListUserProgress :many
SELECT * FROM user_progress
WHERE user_id = $1
ORDER BY last_attempt DESC;

-- name: CreateUserProgress :one
INSERT INTO user_progress (
  user_id, word_id, correct_count, incorrect_count
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUserProgress :one
UPDATE user_progress
SET correct_count = $1, incorrect_count = $2, last_attempt = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteUserProgress :exec
DELETE FROM user_progress
WHERE id = $1;
