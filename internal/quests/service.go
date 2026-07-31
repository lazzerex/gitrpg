package quests

import (
	"context"
	"log/slog"
	"time"

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

// EvaluateAndSave assigns this week's quests on first sight, updates progress
// from stat deltas, and returns quests newly completed by this evaluation.
func (s *Service) EvaluateAndSave(ctx context.Context, userID int64, gs *github.Stats) ([]Quest, error) {
	week := weekStart(time.Now())
	existing, err := s.store.getWeek(ctx, userID, week)
	if err != nil {
		return nil, err
	}

	var completed []Quest
	for _, q := range catalog {
		row, ok := existing[q.Slug]
		if !ok {
			if err := s.store.assign(ctx, userID, week, q.Slug, q.metric(gs), q.XP); err != nil {
				return nil, err
			}
			continue
		}
		if row.Completed {
			continue
		}
		progress := progressOf(q.metric(gs), row.Baseline, q.Target)
		done := progress >= q.Target
		if err := s.store.updateProgress(ctx, userID, week, q.Slug, progress, done); err != nil {
			return nil, err
		}
		if done {
			completed = append(completed, q)
		}
	}
	if len(completed) > 0 {
		s.logger.Info("quests completed", "user_id", userID, "count", len(completed))
	}
	return completed, nil
}

// TotalBonusXP returns the XP sum of all quests the user ever completed.
func (s *Service) TotalBonusXP(ctx context.Context, userID int64) (int, error) {
	return s.store.totalCompletedXP(ctx, userID)
}

// UserQuest is a quest with the user's current-week state, for display.
type UserQuest struct {
	Name      string
	Target    int
	Progress  int
	XP        int
	Completed bool
}

// GetForUser returns the current week's quests in catalog order. Quests not
// yet assigned (user hasn't synced this week) show zero progress.
func (s *Service) GetForUser(ctx context.Context, userID int64) ([]UserQuest, error) {
	rows, err := s.store.getWeek(ctx, userID, weekStart(time.Now()))
	if err != nil {
		return nil, err
	}
	out := make([]UserQuest, 0, len(catalog))
	for _, q := range catalog {
		uq := UserQuest{Name: q.Name, Target: q.Target, XP: q.XP}
		if row, ok := rows[q.Slug]; ok {
			uq.Progress = row.Progress
			uq.Completed = row.Completed
			if row.Completed {
				uq.Progress = q.Target
			}
		}
		out = append(out, uq)
	}
	return out, nil
}
