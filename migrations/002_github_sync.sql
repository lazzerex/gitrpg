-- +goose Up

CREATE TABLE github_stats (
    user_id         BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    commits_count   INT         NOT NULL DEFAULT 0,
    prs_merged      INT         NOT NULL DEFAULT 0,
    issues_closed   INT         NOT NULL DEFAULT 0,
    reviews_count   INT         NOT NULL DEFAULT 0,
    stars_received  INT         NOT NULL DEFAULT 0,
    followers_count INT         NOT NULL DEFAULT 0,
    repos_count     INT         NOT NULL DEFAULT 0,
    qualified_repos INT         NOT NULL DEFAULT 0,
    oss_repos_count INT         NOT NULL DEFAULT 0,
    languages       JSONB       NOT NULL DEFAULT '{}',
    longest_streak  INT         NOT NULL DEFAULT 0,
    current_streak  INT         NOT NULL DEFAULT 0,
    active_days_90  INT         NOT NULL DEFAULT 0,
    synced_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE github_syncs (
    id           BIGSERIAL   PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    status       TEXT        NOT NULL DEFAULT 'running',
    points_used  INT         NOT NULL DEFAULT 0,
    error        TEXT
);

CREATE INDEX idx_github_syncs_user_id ON github_syncs(user_id);
CREATE INDEX idx_github_syncs_status  ON github_syncs(status);

-- +goose Down
DROP TABLE IF EXISTS github_syncs;
DROP TABLE IF EXISTS github_stats;
