package achievements

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/github"
)

type Service struct {
	store  *Store
	logger *slog.Logger
}

func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{store: NewStore(db), logger: logger}
}

// EvaluateAndSave checks which achievements the user earned and persists new ones.
func (s *Service) EvaluateAndSave(ctx context.Context, userID int64, gs *github.Stats) error {
	slugs := Evaluate(gs)
	if err := s.store.Upsert(ctx, userID, slugs); err != nil {
		return err
	}
	s.logger.Debug("achievements evaluated", "user_id", userID, "earned", len(slugs))
	return nil
}

// GetForUser returns the full achievement list with earned status for a user.
func (s *Service) GetForUser(ctx context.Context, userID int64) ([]UserAchievement, error) {
	slugs, err := s.store.GetSlugs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return BuildUserAchievements(slugs), nil
}
