package characters

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/lazzerex/gitrpg/internal/stats"
	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_UpsertAndGet(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "char-store-user")

	if _, err := s.getByUserID(ctx, userID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("before upsert: err = %v, want pgx.ErrNoRows", err)
	}

	char := &stats.Character{
		UserID: userID, TotalXP: 5000, Level: 8, XPIntoLevel: 200, XPForLevel: 900,
		Strength: 40, Intelligence: 55, Wisdom: 20, Dexterity: 33, Charisma: 12,
		Class: "Guardian", Title: "The Maintainer",
	}
	if err := s.upsert(ctx, char); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := s.getByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("getByUserID: %v", err)
	}
	if got.TotalXP != 5000 || got.Level != 8 || got.Class != "Guardian" || got.Title != "The Maintainer" {
		t.Errorf("getByUserID = %+v, want the values written", got)
	}
	if got.UpdatedAt.IsZero() {
		t.Error("UpdatedAt not populated")
	}

	char.Level = 9
	char.TotalXP = 6100
	if err := s.upsert(ctx, char); err != nil {
		t.Fatalf("upsert update: %v", err)
	}
	got, err = s.getByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("getByUserID after update: %v", err)
	}
	if got.Level != 9 || got.TotalXP != 6100 {
		t.Errorf("after conflict update: level %d xp %d, want 9 / 6100", got.Level, got.TotalXP)
	}
}
