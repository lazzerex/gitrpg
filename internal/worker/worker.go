package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/lazzerex/gitrpg/internal/achievements"
	"github.com/lazzerex/gitrpg/internal/characters"
	"github.com/lazzerex/gitrpg/internal/equipment"
	"github.com/lazzerex/gitrpg/internal/events"
	"github.com/lazzerex/gitrpg/internal/github"
	"github.com/lazzerex/gitrpg/internal/stats"
	"github.com/lazzerex/gitrpg/internal/users"
)

// Worker runs background GitHub sync and character computation jobs.
type Worker struct {
	github       *github.Service
	characters   *characters.Service
	achievements *achievements.Service
	equipment    *equipment.Service
	userStore    *users.Store
	bus          *events.Bus
	logger       *slog.Logger
}

// New creates a Worker.
func New(githubSvc *github.Service, charSvc *characters.Service, achSvc *achievements.Service, eqSvc *equipment.Service, userStore *users.Store, bus *events.Bus, logger *slog.Logger) *Worker {
	return &Worker{
		github:       githubSvc,
		characters:   charSvc,
		achievements: achSvc,
		equipment:    eqSvc,
		userStore:    userStore,
		bus:          bus,
		logger:       logger,
	}
}

// Start runs the periodic re-sync loop until ctx is cancelled.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	w.syncAll(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.syncAll(ctx)
		}
	}
}

func (w *Worker) syncAll(ctx context.Context) {
	allUsers, err := w.userStore.ListAll(ctx)
	if err != nil {
		w.logger.Error("worker: list users failed", "error", err)
		return
	}
	w.logger.Info("worker: periodic sync starting", "count", len(allUsers))
	for _, u := range allUsers {
		select {
		case <-ctx.Done():
			return
		default:
		}
		needsSync, err := w.github.NeedsSync(ctx, u)
		if err != nil {
			w.logger.Warn("worker: activity check failed, syncing anyway", "user_id", u.ID, "login", u.Login, "error", err)
		} else if !needsSync {
			w.logger.Info("worker: skip sync, no new activity", "user_id", u.ID, "login", u.Login)
			continue
		}
		syncCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		w.syncOne(syncCtx, u)
		cancel()
	}
}

func (w *Worker) LatestSyncStatus(ctx context.Context, userID int64) (*github.SyncStatus, error) {
	return w.github.LatestSyncStatus(ctx, userID)
}

// SyncUser triggers an immediate background sync and character computation for
// the given user.
func (w *Worker) SyncUser(user *users.User) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		w.syncOne(ctx, user)
	}()
}

func (w *Worker) syncOne(ctx context.Context, user *users.User) {
	if err := w.github.Sync(ctx, user); err != nil {
		w.logger.Error("sync failed", "user_id", user.ID, "login", user.Login, "error", err)
		return
	}
	w.bus.Publish(ctx, events.Event{UserID: user.ID, Type: events.TypeGithubSynced})

	gs, err := w.github.GetStats(ctx, user.ID)
	if err != nil {
		w.logger.Error("get stats failed", "user_id", user.ID, "error", err)
		return
	}

	prev, err := w.characters.GetByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		w.logger.Error("previous character load failed", "user_id", user.ID, "error", err)
	}

	char := stats.Calculate(gs)
	if err := w.characters.Upsert(ctx, char); err != nil {
		w.logger.Error("character upsert failed", "user_id", user.ID, "error", err)
		return
	}
	if prev != nil && char.Level > prev.Level {
		w.bus.Publish(ctx, events.Event{UserID: user.ID, Type: events.TypeLevelUp,
			Payload: map[string]any{"from": prev.Level, "to": char.Level}})
	}

	unlocked, err := w.achievements.EvaluateAndSave(ctx, user.ID, gs)
	if err != nil {
		w.logger.Error("achievement eval failed", "user_id", user.ID, "error", err)
	}
	for _, slug := range unlocked {
		w.bus.Publish(ctx, events.Event{UserID: user.ID, Type: events.TypeAchievementUnlocked,
			Payload: map[string]any{"slug": slug}})
	}

	eqChanged, err := w.equipment.EvaluateAndSave(ctx, user.ID, gs)
	if err != nil {
		w.logger.Error("equipment eval failed", "user_id", user.ID, "error", err)
	} else if eqChanged {
		w.bus.Publish(ctx, events.Event{UserID: user.ID, Type: events.TypeEquipmentChanged})
	}

	w.logger.Info("sync complete",
		"user_id", user.ID,
		"login", user.Login,
		"level", char.Level,
		"xp", char.TotalXP,
		"class", char.Class,
	)
}
