package leaderboards

import (
	"context"
	"testing"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_XPPageOrdersByTotalXP(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()

	before, err := s.countCharacters(ctx)
	if err != nil {
		t.Fatalf("countCharacters: %v", err)
	}

	xp := map[string]int{"lb-low": 100, "lb-mid": 5000, "lb-high": 90000}
	for login, total := range xp {
		userID := testdb.User(t, pool, login)
		if _, err := pool.Exec(ctx, `
			INSERT INTO characters (user_id, total_xp, level, class, title)
			VALUES ($1, $2, 1, 'Guardian', 'The Adventurer')`, userID, total); err != nil {
			t.Fatalf("insert character: %v", err)
		}
	}

	after, err := s.countCharacters(ctx)
	if err != nil {
		t.Fatalf("countCharacters after insert: %v", err)
	}
	if after != before+3 {
		t.Errorf("countCharacters = %d, want %d", after, before+3)
	}

	entries, err := s.getXPPage(ctx, 100, 0)
	if err != nil {
		t.Fatalf("getXPPage: %v", err)
	}
	for i := 1; i < len(entries); i++ {
		if entries[i-1].TotalXP < entries[i].TotalXP {
			t.Fatalf("entries not ordered by XP desc at %d: %d then %d",
				i, entries[i-1].TotalXP, entries[i].TotalXP)
		}
		if entries[i-1].Rank > entries[i].Rank {
			t.Fatalf("ranks not ascending at %d: %d then %d", i, entries[i-1].Rank, entries[i].Rank)
		}
	}

	if len(entries) > 1 {
		second, err := s.getXPPage(ctx, 1, 1)
		if err != nil {
			t.Fatalf("getXPPage offset: %v", err)
		}
		if len(second) != 1 {
			t.Fatalf("limit 1 returned %d rows", len(second))
		}
		if second[0].Login == entries[0].Login {
			t.Error("offset 1 returned the first row again")
		}
	}
}
