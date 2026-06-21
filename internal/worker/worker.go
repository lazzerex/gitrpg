package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/lazzerex/gitrpg/internal/characters"
	"github.com/lazzerex/gitrpg/internal/github"
	"github.com/lazzerex/gitrpg/internal/stats"
	"github.com/lazzerex/gitrpg/internal/users"
)

// Worker runs background GitHub sync and character computation jobs.
type Worker struct {
	github     *github.Service
	characters *characters.Service
	logger     *slog.Logger
}

// New creates a Worker.
func New(githubSvc *github.Service, charSvc *characters.Service, logger *slog.Logger) *Worker {
	return &Worker{github: githubSvc, characters: charSvc, logger: logger}
}

// Start runs the periodic re-sync loop until ctx is cancelled.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// TODO: list all users and re-sync each (users.Store.ListAll not yet implemented)
			w.logger.Info("worker: periodic sync tick")
		}
	}
}

// SyncUser triggers an immediate sync and character computation for the given user.
func (w *Worker) SyncUser(user *users.User) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if err := w.github.Sync(ctx, user); err != nil {
			w.logger.Error("sync failed", "user_id", user.ID, "login", user.Login, "error", err)
			return
		}

		gs, err := w.github.GetStats(ctx, user.ID)
		if err != nil {
			w.logger.Error("get stats failed", "user_id", user.ID, "error", err)
			return
		}

		char := stats.Calculate(gs)
		if err := w.characters.Upsert(ctx, char); err != nil {
			w.logger.Error("character upsert failed", "user_id", user.ID, "error", err)
			return
		}

		w.logger.Info("sync complete",
			"user_id", user.ID,
			"login", user.Login,
			"level", char.Level,
			"xp", char.TotalXP,
			"class", char.Class,
		)
	}()
}
