package achievements

import (
	"context"
	"testing"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_UpsertAndGetSlugs(t *testing.T) {
	pool := testdb.Pool(t)
	s := NewStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "ach-store-user")

	slugs, err := s.GetSlugs(ctx, userID)
	if err != nil || len(slugs) != 0 {
		t.Fatalf("empty user: slugs = %v, err = %v", slugs, err)
	}

	if err := s.Upsert(ctx, userID, nil); err != nil {
		t.Fatalf("Upsert with no slugs: %v", err)
	}

	if err := s.Upsert(ctx, userID, []string{"first-commit", "first-pr"}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	slugs, err = s.GetSlugs(ctx, userID)
	if err != nil {
		t.Fatalf("GetSlugs: %v", err)
	}
	if len(slugs) != 2 {
		t.Fatalf("slugs = %v, want 2", slugs)
	}

	if err := s.Upsert(ctx, userID, []string{"first-commit", "first-repo"}); err != nil {
		t.Fatalf("Upsert overlapping: %v", err)
	}
	slugs, err = s.GetSlugs(ctx, userID)
	if err != nil {
		t.Fatalf("GetSlugs after overlap: %v", err)
	}
	if len(slugs) != 3 {
		t.Errorf("slugs = %v, want 3 — re-awarding must not duplicate", slugs)
	}
}
