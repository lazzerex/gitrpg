package characters

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/stats"
)

// Service manages character persistence.
type Service struct {
	store  *store
	logger *slog.Logger
}

// NewService creates a Service backed by the given DB pool.
func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{store: newStore(db), logger: logger}
}

// Upsert persists a computed character.
func (s *Service) Upsert(ctx context.Context, c *stats.Character) error {
	return s.store.upsert(ctx, c)
}

// GetByUserID returns the stored character for the given user.
func (s *Service) GetByUserID(ctx context.Context, userID int64) (*stats.Character, error) {
	return s.store.getByUserID(ctx, userID)
}
