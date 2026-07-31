-- +goose Up
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    github_id   BIGINT      NOT NULL UNIQUE,
    login       TEXT        NOT NULL,
    name        TEXT,
    avatar_url  TEXT,
    email       TEXT,
    access_token TEXT       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_login ON users(login);

-- +goose Down
DROP TABLE IF EXISTS users;
