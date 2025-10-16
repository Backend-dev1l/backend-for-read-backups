-- name: GetUserWordSet :one
SELECT * FROM user_word_sets
WHERE id = $1 LIMIT 1;

-- name: ListUserWordSets :many
SELECT * FROM user_word_sets
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateUserWordSet :one
INSERT INTO user_word_sets (
  user_id, word_set_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteUserWordSet :exec
DELETE FROM user_word_sets
WHERE id = $1;
