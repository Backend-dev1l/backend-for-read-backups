-- +goose Up
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    word_id UUID NOT NULL,
    correct_count INTEGER NOT NULL DEFAULT 0 CHECK (correct_count >= 0),
    incorrect_count INTEGER NOT NULL DEFAULT 0 CHECK (incorrect_count >= 0),
    last_attempt TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, word_id)
);

-- +goose Down
DROP TABLE IF NOT EXISTS user_progress;