-- name: GetUserWordSet :one
SELECT * FROM user_word_sets
WHERE id = $1 LIMIT 1;

-- name: ListUserWordSets :many
SELECT * FROM user_word_sets
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3; 

-- name: CreateUserWordSet :one
INSERT INTO user_word_sets (
  user_id, word_set_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUserWordSet :one
UPDATE user_word_sets
SET word_set_id = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserWordSet :exec
DELETE FROM user_word_sets
WHERE id = $1;
