package server

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func newRedisServer(t *testing.T) *Server {
	t.Helper()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatalf("parse REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skipf("redis unreachable: %v", err)
	}
	t.Cleanup(func() { _ = rdb.Close() })
	return &Server{redis: rdb, logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func TestSyncStartMarker_RoundTrip(t *testing.T) {
	s := newRedisServer(t)
	ctx := context.Background()
	const userID = -9001
	t.Cleanup(func() { s.clearSyncStart(ctx, userID) })

	if _, ok, err := s.syncStart(ctx, userID); err != nil || ok {
		t.Fatalf("before set: ok = %v, err = %v, want false/nil", ok, err)
	}

	want := time.Now().Truncate(time.Millisecond)
	if err := s.setSyncStart(ctx, userID, want); err != nil {
		t.Fatalf("setSyncStart: %v", err)
	}

	got, ok, err := s.syncStart(ctx, userID)
	if err != nil || !ok {
		t.Fatalf("after set: ok = %v, err = %v, want true/nil", ok, err)
	}
	if !got.Equal(want) {
		t.Errorf("syncStart = %v, want %v", got, want)
	}

	s.clearSyncStart(ctx, userID)
	if _, ok, _ := s.syncStart(ctx, userID); ok {
		t.Error("marker still present after clear")
	}
}

func TestSyncStartMarker_VisibleToAnotherServer(t *testing.T) {
	writer := newRedisServer(t)
	reader := newRedisServer(t)
	ctx := context.Background()
	const userID = -9002
	t.Cleanup(func() { writer.clearSyncStart(ctx, userID) })

	want := time.Now().Truncate(time.Millisecond)
	if err := writer.setSyncStart(ctx, userID, want); err != nil {
		t.Fatalf("setSyncStart: %v", err)
	}
	got, ok, err := reader.syncStart(ctx, userID)
	if err != nil || !ok {
		t.Fatalf("second instance: ok = %v, err = %v, want true/nil", ok, err)
	}
	if !got.Equal(want) {
		t.Errorf("syncStart = %v, want %v", got, want)
	}
}
