package github

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/users"
)

// Service orchestrates GitHub data sync for a user.
type Service struct {
	store  *store
	logger *slog.Logger
}

// NewService creates a Service backed by the given database pool.
func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{store: newStore(db), logger: logger}
}

// Sync fetches GitHub data for user and persists it.
func (s *Service) Sync(ctx context.Context, user *users.User) error {
	syncID, err := s.store.startSync(ctx, user.ID)
	if err != nil {
		return err
	}

	raw, fetchErr := fetch(ctx, user.AccessToken, user.Login, s.logger)
	if fetchErr != nil {
		_ = s.store.completeSync(ctx, syncID, 0, fetchErr)
		return fetchErr
	}

	stats := process(user.ID, raw)

	if err := s.store.upsertStats(ctx, stats); err != nil {
		_ = s.store.completeSync(ctx, syncID, raw.PointsUsed, err)
		return err
	}

	return s.store.completeSync(ctx, syncID, raw.PointsUsed, nil)
}

// GetStats returns the most recently synced stats for a user.
func (s *Service) GetStats(ctx context.Context, userID int64) (*Stats, error) {
	return s.store.getStats(ctx, userID)
}
