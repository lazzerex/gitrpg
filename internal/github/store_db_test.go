package github

import (
	"context"
	"errors"
	"testing"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_UpsertAndGetStats(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "gh-stats-user")

	want := &Stats{
		UserID: userID, CommitsCount: 120, PRsMerged: 7, IssuesClosed: 3, ReviewsCount: 11,
		StarsReceived: 42, FollowersCount: 9, ReposCount: 5, QualifiedRepos: 4, OSSReposCount: 2,
		Languages: map[string]int{"Go": 90, "C++": 30}, LongestStreak: 15, CurrentStreak: 4, ActiveDays90: 33,
	}
	if err := s.upsertStats(ctx, want); err != nil {
		t.Fatalf("upsertStats: %v", err)
	}

	got, err := s.getStats(ctx, userID)
	if err != nil {
		t.Fatalf("getStats: %v", err)
	}
	if got.CommitsCount != 120 || got.QualifiedRepos != 4 || got.LongestStreak != 15 {
		t.Errorf("getStats = %+v, want the values written", got)
	}
	if got.Languages["Go"] != 90 || got.Languages["C++"] != 30 {
		t.Errorf("Languages = %v, want Go:90 C++:30 through jsonb", got.Languages)
	}

	want.CommitsCount = 200
	want.Languages = map[string]int{"Rust": 5}
	if err := s.upsertStats(ctx, want); err != nil {
		t.Fatalf("upsertStats update: %v", err)
	}
	got, err = s.getStats(ctx, userID)
	if err != nil {
		t.Fatalf("getStats after update: %v", err)
	}
	if got.CommitsCount != 200 {
		t.Errorf("CommitsCount = %d, want 200 after conflict update", got.CommitsCount)
	}
	if _, stale := got.Languages["Go"]; stale {
		t.Errorf("Languages = %v, want the old map replaced", got.Languages)
	}
}

func TestStore_GetLastSyncedAt(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "gh-lastsync-user")

	if _, ok, err := s.getLastSyncedAt(ctx, userID); err != nil || ok {
		t.Fatalf("before stats: ok = %v, err = %v, want false/nil", ok, err)
	}
	if err := s.upsertStats(ctx, &Stats{UserID: userID, Languages: map[string]int{}}); err != nil {
		t.Fatalf("upsertStats: %v", err)
	}
	if _, ok, err := s.getLastSyncedAt(ctx, userID); err != nil || !ok {
		t.Fatalf("after stats: ok = %v, err = %v, want true/nil", ok, err)
	}
}

func TestStore_SyncLifecycle(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "gh-sync-user")

	if _, _, _, ok, err := s.getLatestSync(ctx, userID); err != nil || ok {
		t.Fatalf("no syncs yet: ok = %v, err = %v, want false/nil", ok, err)
	}

	cases := []struct {
		name       string
		syncErr    error
		wantStatus string
		wantErrMsg bool
	}{
		{"success", nil, "success", false},
		{"failure", errors.New("boom"), "failed", true},
		{"rejected token", statusError(401), StatusUnauthorized, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := s.startSync(ctx, userID)
			if err != nil {
				t.Fatalf("startSync: %v", err)
			}
			if err := s.completeSync(ctx, id, 17, c.syncErr); err != nil {
				t.Fatalf("completeSync: %v", err)
			}

			status, completedAt, errMsg, ok, err := s.getLatestSync(ctx, userID)
			if err != nil || !ok {
				t.Fatalf("getLatestSync: ok = %v, err = %v", ok, err)
			}
			if status != c.wantStatus {
				t.Errorf("status = %q, want %q", status, c.wantStatus)
			}
			if completedAt == nil {
				t.Error("completedAt is nil after completeSync")
			}
			if (errMsg != "") != c.wantErrMsg {
				t.Errorf("errMsg = %q, want non-empty = %v", errMsg, c.wantErrMsg)
			}
		})
	}
}
