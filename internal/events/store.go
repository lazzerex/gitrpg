package events

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

func (s *store) insert(ctx context.Context, e Event) error {
	payload := e.Payload
	if payload == nil {
		payload = map[string]any{}
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(ctx,
		`INSERT INTO events (user_id, type, payload) VALUES ($1, $2, $3::jsonb)`,
		e.UserID, e.Type, b)
	return err
}
