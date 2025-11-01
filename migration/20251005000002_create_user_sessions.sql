
-- +goose Up
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at TIMESTAMPTZ NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'completed', 'abandoned'))
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_status ON user_sessions (status) WHERE status = 'active';

-- +goose Down
DROP TABLE IF EXISTS user_sessions;
DROP INDEX IF EXISTS idx_user_sessions_user_id
DROP INDEX IF EXISTS idx_user_sessions_status 


