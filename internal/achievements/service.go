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

// EvaluateAndSave checks which achievements the user earned, persists new ones,
// and returns the newly unlocked slugs.
func (s *Service) EvaluateAndSave(ctx context.Context, userID int64, gs *github.Stats) ([]string, error) {
	earned := Evaluate(gs)
	existing, err := s.store.GetSlugs(ctx, userID)
	if err != nil {
		return nil, err
	}
	unlocked := diffSlugs(earned, existing)
	if err := s.store.Upsert(ctx, userID, unlocked); err != nil {
		return nil, err
	}
	s.logger.Debug("achievements evaluated", "user_id", userID, "earned", len(earned), "new", len(unlocked))
	return unlocked, nil
}

func diffSlugs(earned, existing []string) []string {
	have := make(map[string]struct{}, len(existing))
	for _, s := range existing {
		have[s] = struct{}{}
	}
	var out []string
	for _, s := range earned {
		if _, ok := have[s]; !ok {
			out = append(out, s)
		}
	}
	return out
}

// GetForUser returns the full achievement list with earned status for a user.
func (s *Service) GetForUser(ctx context.Context, userID int64) ([]UserAchievement, error) {
	slugs, err := s.store.GetSlugs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return BuildUserAchievements(slugs), nil
}
