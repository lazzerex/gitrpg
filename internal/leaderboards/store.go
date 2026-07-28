package leaderboards

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

func (s *store) countCharacters(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM characters`).Scan(&count)
	return count, err
}

func (s *store) getXPPage(ctx context.Context, limit, offset int) ([]Entry, error) {
	rows, err := s.db.Query(ctx, `
		SELECT RANK() OVER (ORDER BY c.total_xp DESC) AS rank,
		       u.login, u.avatar_url, c.class, c.level, c.total_xp
		FROM characters c
		JOIN users u ON u.id = c.user_id
		ORDER BY c.total_xp DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var avatarURL *string
		if err := rows.Scan(&e.Rank, &e.Login, &avatarURL, &e.Class, &e.Level, &e.TotalXP); err != nil {
			return nil, err
		}
		if avatarURL != nil {
			e.AvatarURL = *avatarURL
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
