-- name: GetUserSession :one
SELECT * FROM user_sessions
WHERE id = $1 LIMIT 1;

-- name: ListUserSessions :many
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY started_at DESC;

-- name: ListActiveSessions :many
SELECT * FROM user_sessions
WHERE status = 'active'
ORDER BY started_at DESC;

-- name: CreateUserSession :one
INSERT INTO user_sessions (
  user_id, status
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET status = $1, ended_at = $2
WHERE id = $3
RETURNING *;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions
WHERE id = $1;
