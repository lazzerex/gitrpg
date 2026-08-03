package quests

import (
	"context"
	"testing"
	"time"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_AssignAndProgress(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "quest-progress-user")
	week := weekStart(time.Now())

	if err := s.assign(ctx, userID, week, "weekly-prs", 10, 150); err != nil {
		t.Fatalf("assign: %v", err)
	}
	rows, err := s.getWeek(ctx, userID, week)
	if err != nil {
		t.Fatalf("getWeek: %v", err)
	}
	row, ok := rows["weekly-prs"]
	if !ok {
		t.Fatal("assigned quest missing from getWeek")
	}
	if row.Baseline != 10 || row.Progress != 0 || row.Completed {
		t.Errorf("row = %+v, want baseline 10, progress 0, not completed", row)
	}

	if err := s.assign(ctx, userID, week, "weekly-prs", 999, 999); err != nil {
		t.Fatalf("re-assign: %v", err)
	}
	rows, _ = s.getWeek(ctx, userID, week)
	if rows["weekly-prs"].Baseline != 10 {
		t.Errorf("baseline = %d, want 10 — re-assign must not overwrite", rows["weekly-prs"].Baseline)
	}

	if err := s.updateProgress(ctx, userID, week, "weekly-prs", 1, false); err != nil {
		t.Fatalf("updateProgress: %v", err)
	}
	rows, _ = s.getWeek(ctx, userID, week)
	if rows["weekly-prs"].Progress != 1 || rows["weekly-prs"].Completed {
		t.Errorf("row = %+v, want progress 1 and not completed", rows["weekly-prs"])
	}
}

func TestStore_CompletionIsSticky(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "quest-complete-user")
	week := weekStart(time.Now())

	if err := s.assign(ctx, userID, week, "weekly-commits", 0, 100); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if err := s.updateProgress(ctx, userID, week, "weekly-commits", 20, true); err != nil {
		t.Fatalf("complete: %v", err)
	}

	var firstCompletedAt time.Time
	if err := pool.QueryRow(ctx,
		`SELECT completed_at FROM quests WHERE user_id=$1 AND slug='weekly-commits' AND week_start=$2`,
		userID, week).Scan(&firstCompletedAt); err != nil {
		t.Fatalf("read completed_at: %v", err)
	}

	if err := s.updateProgress(ctx, userID, week, "weekly-commits", 20, true); err != nil {
		t.Fatalf("second complete: %v", err)
	}
	var secondCompletedAt time.Time
	if err := pool.QueryRow(ctx,
		`SELECT completed_at FROM quests WHERE user_id=$1 AND slug='weekly-commits' AND week_start=$2`,
		userID, week).Scan(&secondCompletedAt); err != nil {
		t.Fatalf("re-read completed_at: %v", err)
	}
	if !firstCompletedAt.Equal(secondCompletedAt) {
		t.Errorf("completed_at moved from %v to %v; completion must be recorded once",
			firstCompletedAt, secondCompletedAt)
	}
}

func TestStore_TotalCompletedXP(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "quest-xp-user")
	week := weekStart(time.Now())

	total, err := s.totalCompletedXP(ctx, userID)
	if err != nil || total != 0 {
		t.Fatalf("empty total = %d, err = %v, want 0/nil", total, err)
	}

	if err := s.assign(ctx, userID, week, "weekly-prs", 0, 150); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if err := s.assign(ctx, userID, week, "weekly-reviews", 0, 100); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if err := s.updateProgress(ctx, userID, week, "weekly-prs", 2, true); err != nil {
		t.Fatalf("complete: %v", err)
	}

	total, err = s.totalCompletedXP(ctx, userID)
	if err != nil {
		t.Fatalf("totalCompletedXP: %v", err)
	}
	if total != 150 {
		t.Errorf("total = %d, want 150 — incomplete quests must not count", total)
	}
}
