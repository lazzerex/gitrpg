package quests

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

type questRow struct {
	Slug      string
	Baseline  int
	Progress  int
	XP        int
	Completed bool
}

func (s *store) getWeek(ctx context.Context, userID int64, week time.Time) (map[string]questRow, error) {
	rows, err := s.db.Query(ctx, `
		SELECT slug, baseline, progress, xp, completed_at IS NOT NULL
		FROM quests WHERE user_id = $1 AND week_start = $2`, userID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]questRow)
	for rows.Next() {
		var r questRow
		if err := rows.Scan(&r.Slug, &r.Baseline, &r.Progress, &r.XP, &r.Completed); err != nil {
			return nil, err
		}
		out[r.Slug] = r
	}
	return out, rows.Err()
}

func (s *store) assign(ctx context.Context, userID int64, week time.Time, slug string, baseline, xp int) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO quests (user_id, slug, week_start, baseline, xp)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING`, userID, slug, week, baseline, xp)
	return err
}

func (s *store) updateProgress(ctx context.Context, userID int64, week time.Time, slug string, progress int, completed bool) error {
	if completed {
		_, err := s.db.Exec(ctx, `
			UPDATE quests SET progress = $4, completed_at = now()
			WHERE user_id = $1 AND week_start = $2 AND slug = $3 AND completed_at IS NULL`,
			userID, week, slug, progress)
		return err
	}
	_, err := s.db.Exec(ctx, `
		UPDATE quests SET progress = $4
		WHERE user_id = $1 AND week_start = $2 AND slug = $3`,
		userID, week, slug, progress)
	return err
}

func (s *store) totalCompletedXP(ctx context.Context, userID int64) (int, error) {
	var total int
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(xp), 0) FROM quests
		WHERE user_id = $1 AND completed_at IS NOT NULL`, userID).Scan(&total)
	return total, err
}
