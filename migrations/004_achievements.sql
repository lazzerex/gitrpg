-- +goose Up
CREATE TABLE user_achievements (
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug        TEXT NOT NULL,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, slug)
);

CREATE INDEX user_achievements_user_id ON user_achievements(user_id);

-- +goose Down
DROP TABLE IF EXISTS user_achievements;
