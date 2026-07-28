package equipment

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/github"
)

type Service struct {
	store  *store
	logger *slog.Logger
}

func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{store: newStore(db), logger: logger}
}

func (s *Service) EvaluateAndSave(ctx context.Context, userID int64, gs *github.Stats) error {
	loadout := Evaluate(gs)
	if err := s.store.upsert(ctx, userID, loadout); err != nil {
		return err
	}
	s.logger.Debug("equipment evaluated", "user_id", userID,
		"weapon", slugOf(loadout.Weapon), "shield", slugOf(loadout.Shield), "accessory", slugOf(loadout.Accessory))
	return nil
}

// Unsynced users get a zero-value Loadout, not an error.
func (s *Service) GetForUser(ctx context.Context, userID int64) (Loadout, error) {
	return s.store.getByUserID(ctx, userID)
}

func slugOf(it *Item) string {
	if it == nil {
		return ""
	}
	return it.Slug
}
