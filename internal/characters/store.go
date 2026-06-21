package characters

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/stats"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

func (s *store) upsert(ctx context.Context, c *stats.Character) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO characters (
			user_id, total_xp, level, xp_into_level, xp_for_level,
			strength, intelligence, wisdom, dexterity, charisma,
			class, title, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,now())
		ON CONFLICT (user_id) DO UPDATE SET
			total_xp      = EXCLUDED.total_xp,
			level         = EXCLUDED.level,
			xp_into_level = EXCLUDED.xp_into_level,
			xp_for_level  = EXCLUDED.xp_for_level,
			strength      = EXCLUDED.strength,
			intelligence  = EXCLUDED.intelligence,
			wisdom        = EXCLUDED.wisdom,
			dexterity     = EXCLUDED.dexterity,
			charisma      = EXCLUDED.charisma,
			class         = EXCLUDED.class,
			title         = EXCLUDED.title,
			updated_at    = now()`,
		c.UserID, c.TotalXP, c.Level, c.XPIntoLevel, c.XPForLevel,
		c.Strength, c.Intelligence, c.Wisdom, c.Dexterity, c.Charisma,
		c.Class, c.Title,
	)
	return err
}

func (s *store) getByUserID(ctx context.Context, userID int64) (*stats.Character, error) {
	c := &stats.Character{UserID: userID}
	err := s.db.QueryRow(ctx, `
		SELECT total_xp, level, xp_into_level, xp_for_level,
		       strength, intelligence, wisdom, dexterity, charisma,
		       class, title
		FROM characters WHERE user_id = $1`, userID).
		Scan(&c.TotalXP, &c.Level, &c.XPIntoLevel, &c.XPForLevel,
			&c.Strength, &c.Intelligence, &c.Wisdom, &c.Dexterity, &c.Charisma,
			&c.Class, &c.Title)
	if err != nil {
		return nil, err
	}
	return c, nil
}
