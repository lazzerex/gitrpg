-- +goose Up
CREATE TABLE equipment (
    user_id        BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    weapon_slug    TEXT,
    shield_slug    TEXT,
    accessory_slug TEXT,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS equipment;
