package equipment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

func (s *store) upsert(ctx context.Context, userID int64, l Loadout) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO equipment (user_id, weapon_slug, shield_slug, accessory_slug, updated_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (user_id) DO UPDATE SET
			weapon_slug    = EXCLUDED.weapon_slug,
			shield_slug    = EXCLUDED.shield_slug,
			accessory_slug = EXCLUDED.accessory_slug,
			updated_at     = now()`,
		userID, slugPtr(l.Weapon), slugPtr(l.Shield), slugPtr(l.Accessory),
	)
	return err
}

func (s *store) getByUserID(ctx context.Context, userID int64) (Loadout, error) {
	var weapon, shield, accessory *string
	err := s.db.QueryRow(ctx, `
		SELECT weapon_slug, shield_slug, accessory_slug
		FROM equipment WHERE user_id = $1`, userID).
		Scan(&weapon, &shield, &accessory)
	if errors.Is(err, pgx.ErrNoRows) {
		return Loadout{}, nil
	}
	if err != nil {
		return Loadout{}, err
	}
	return Loadout{
		Weapon:    itemFromSlug(weapon),
		Shield:    itemFromSlug(shield),
		Accessory: itemFromSlug(accessory),
	}, nil
}

func slugPtr(it *Item) *string {
	if it == nil {
		return nil
	}
	return &it.Slug
}

func itemFromSlug(slug *string) *Item {
	if slug == nil {
		return nil
	}
	item, ok := bySlug[*slug]
	if !ok {
		return nil
	}
	return &item
}
