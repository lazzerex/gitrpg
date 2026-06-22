package achievements

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// Upsert awards achievement slugs to a user, skipping any already earned.
func (s *Store) Upsert(ctx context.Context, userID int64, slugs []string) error {
	if len(slugs) == 0 {
		return nil
	}
	const q = `
		INSERT INTO user_achievements (user_id, slug)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	for _, slug := range slugs {
		if _, err := s.db.Exec(ctx, q, userID, slug); err != nil {
			return err
		}
	}
	return nil
}

// GetSlugs returns all achievement slugs earned by a user.
func (s *Store) GetSlugs(ctx context.Context, userID int64) ([]string, error) {
	const q = `SELECT slug FROM user_achievements WHERE user_id = $1 ORDER BY unlocked_at`
	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, rows.Err()
}
