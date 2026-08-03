package equipment

import (
	"context"
	"testing"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func TestStore_UpsertAndGet(t *testing.T) {
	pool := testdb.Pool(t)
	s := newStore(pool)
	ctx := context.Background()
	userID := testdb.User(t, pool, "equip-store-user")

	empty, err := s.getByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("getByUserID with no row: %v", err)
	}
	if empty.Any() {
		t.Errorf("missing row returned %+v, want the zero Loadout", empty)
	}

	weapon, _ := Lookup("go-compiler")
	shield, _ := Lookup("csharp-bulwark")
	if err := s.upsert(ctx, userID, Loadout{Weapon: &weapon, Shield: &shield}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := s.getByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("getByUserID: %v", err)
	}
	if got.Weapon == nil || got.Weapon.Slug != "go-compiler" {
		t.Errorf("Weapon = %+v, want go-compiler", got.Weapon)
	}
	if got.Shield == nil || got.Shield.Slug != "csharp-bulwark" {
		t.Errorf("Shield = %+v, want csharp-bulwark", got.Shield)
	}
	if got.Accessory != nil {
		t.Errorf("Accessory = %+v, want nil for the NULL column", got.Accessory)
	}

	if err := s.upsert(ctx, userID, Loadout{}); err != nil {
		t.Fatalf("upsert cleared loadout: %v", err)
	}
	got, err = s.getByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("getByUserID after clear: %v", err)
	}
	if got.Any() {
		t.Errorf("loadout = %+v, want every slot cleared", got)
	}
}
