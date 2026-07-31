-- +goose Up
CREATE TABLE characters (
    user_id       BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_xp      INT         NOT NULL DEFAULT 0,
    level         INT         NOT NULL DEFAULT 0,
    xp_into_level INT         NOT NULL DEFAULT 0,
    xp_for_level  INT         NOT NULL DEFAULT 100,
    strength      INT         NOT NULL DEFAULT 0,
    intelligence  INT         NOT NULL DEFAULT 0,
    wisdom        INT         NOT NULL DEFAULT 5,
    dexterity     INT         NOT NULL DEFAULT 0,
    charisma      INT         NOT NULL DEFAULT 0,
    class         TEXT        NOT NULL DEFAULT 'Wanderer',
    title         TEXT        NOT NULL DEFAULT 'The Adventurer',
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS characters;
