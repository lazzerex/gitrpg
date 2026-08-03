package events

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TypeGithubSynced        = "github.synced"
	TypeCharacterUpdated    = "character.updated"
	TypeLevelUp             = "level.up"
	TypeAchievementUnlocked = "achievement.unlocked"
	TypeEquipmentChanged    = "equipment.changed"
	TypeQuestCompleted      = "quest.completed"
)

// Event is a domain event tied to a user.
type Event struct {
	UserID  int64
	Type    string
	Payload map[string]any
}

type Handler func(ctx context.Context, e Event)

// Bus persists domain events and dispatches them to in-process subscribers.
type Bus struct {
	store    *store
	logger   *slog.Logger
	handlers map[string][]Handler
}

func NewBus(db *pgxpool.Pool, logger *slog.Logger) *Bus {
	return &Bus{store: newStore(db), logger: logger, handlers: make(map[string][]Handler)}
}

// Subscribe registers a handler for an event type. Register during startup;
// not safe for concurrent use with Publish.
func (b *Bus) Subscribe(eventType string, h Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], h)
}

// Publish persists the event and invokes subscribers synchronously.
// Persistence failure is logged, not returned, so emitting never fails the caller.
func (b *Bus) Publish(ctx context.Context, e Event) {
	if err := b.store.insert(ctx, e); err != nil {
		b.logger.Error("event persist failed", "type", e.Type, "user_id", e.UserID, "error", err)
	}
	b.dispatch(ctx, e)
}

func (b *Bus) dispatch(ctx context.Context, e Event) {
	for _, h := range b.handlers[e.Type] {
		h(ctx, e)
	}
}
