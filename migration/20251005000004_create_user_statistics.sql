-- +goose Up
CREATE TABLE user_statistics (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_words_learned INTEGER NOT NULL DEFAULT 0 CHECK (total_words_learned >= 0),
    accuracy NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (accuracy BETWEEN 0 AND 100),
    total_time INTEGER NOT NULL DEFAULT 0 CHECK (total_time >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF NOT EXISTS user_statistics;