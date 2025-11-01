-- +goose Up
CREATE TABLE user_word_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    word_set_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_user_word_set UNIQUE (user_id, word_set_id)
);

CREATE INDEX IF NOT EXISTS idx_user_word_sets_user_id ON user_word_sets(user_id);
CREATE INDEX IF NOT EXISTS idx_user_word_sets_word_set_id ON user_word_sets(word_set_id);

-- +goose Down
DROP TABLE IF EXISTS user_word_sets;
DROP INDEX IF EXISTS idx_user_word_sets_user_id ON user_word_sets(user_id);
DROP INDEX IF EXISTS idx_user_word_sets_word_set_id ON user_word_sets(word_set_id);






