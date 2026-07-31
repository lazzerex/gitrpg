-- +goose Up
CREATE TABLE quests (
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug         TEXT        NOT NULL,
    week_start   DATE        NOT NULL,
    baseline     INT         NOT NULL,
    progress     INT         NOT NULL DEFAULT 0,
    xp           INT         NOT NULL,
    completed_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, slug, week_start)
);

CREATE INDEX idx_quests_user_completed ON quests (user_id) WHERE completed_at IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS quests;
