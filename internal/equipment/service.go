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

// EvaluateAndSave recomputes the loadout and reports whether it changed.
func (s *Service) EvaluateAndSave(ctx context.Context, userID int64, gs *github.Stats) (bool, error) {
	loadout := Evaluate(gs)
	prev, err := s.store.getByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	if err := s.store.upsert(ctx, userID, loadout); err != nil {
		return false, err
	}
	s.logger.Debug("equipment evaluated", "user_id", userID,
		"weapon", slugOf(loadout.Weapon), "shield", slugOf(loadout.Shield), "accessory", slugOf(loadout.Accessory))
	return !sameLoadout(loadout, prev), nil
}

func sameLoadout(a, b Loadout) bool {
	return slugOf(a.Weapon) == slugOf(b.Weapon) &&
		slugOf(a.Shield) == slugOf(b.Shield) &&
		slugOf(a.Accessory) == slugOf(b.Accessory)
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
