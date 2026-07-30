package github

import (
	"context"
	"log/slog"
	"time"

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

type SyncStatus struct {
	Status      string
	CompletedAt *time.Time
	Error       string
}

// LatestSyncStatus returns nil if the user has never had a sync attempt.
func (s *Service) LatestSyncStatus(ctx context.Context, userID int64) (*SyncStatus, error) {
	status, completedAt, errMsg, ok, err := s.store.getLatestSync(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &SyncStatus{Status: status, CompletedAt: completedAt, Error: errMsg}, nil
}

// maxSyncAge forces a full sync even without new activity: streaks decay and
// stars/followers change while a user is idle.
const maxSyncAge = 24 * time.Hour

// NeedsSync reports whether user has new GitHub activity since their last sync,
// so periodic re-sync can skip idle users instead of burning GraphQL points on
// a full fetch. Users never synced, or synced longer than maxSyncAge ago,
// always need a sync.
func (s *Service) NeedsSync(ctx context.Context, user *users.User) (bool, error) {
	since, ok, err := s.store.getLastSyncedAt(ctx, user.ID)
	if err != nil {
		return false, err
	}
	if !ok || time.Since(since) >= maxSyncAge {
		return true, nil
	}
	return hasActivity(ctx, user.AccessToken, user.Login, since, s.logger)
}
