package server

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/lazzerex/gitrpg/internal/events"
)

func TestInvalidateCardCache(t *testing.T) {
	s := newRedisServer(t)
	ctx := context.Background()
	const login = "test-invalidate-user"
	key := cardCacheKey(login)
	t.Cleanup(func() { s.redis.Del(ctx, key) })

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := InvalidateCardCache(s.redis, logger)

	if err := s.redis.Set(ctx, key, "<svg/>", time.Minute).Err(); err != nil {
		t.Fatalf("seed cache: %v", err)
	}
	handler(ctx, events.Event{UserID: 1, Type: events.TypeCharacterUpdated,
		Payload: map[string]any{"login": login}})

	if n, err := s.redis.Exists(ctx, key).Result(); err != nil || n != 0 {
		t.Errorf("cached card still present: exists = %d, err = %v", n, err)
	}
}

func TestInvalidateCardCache_MissingLoginIsNoop(t *testing.T) {
	s := newRedisServer(t)
	ctx := context.Background()
	const login = "test-noop-user"
	key := cardCacheKey(login)
	t.Cleanup(func() { s.redis.Del(ctx, key) })

	if err := s.redis.Set(ctx, key, "<svg/>", time.Minute).Err(); err != nil {
		t.Fatalf("seed cache: %v", err)
	}
	handler := InvalidateCardCache(s.redis, slog.New(slog.NewTextHandler(io.Discard, nil)))
	handler(ctx, events.Event{UserID: 1, Type: events.TypeCharacterUpdated})

	if n, _ := s.redis.Exists(ctx, key).Result(); n != 1 {
		t.Error("handler with no login in payload must not touch the cache")
	}
}
