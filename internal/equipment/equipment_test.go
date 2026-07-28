package equipment

import (
	"testing"

	"github.com/lazzerex/gitrpg/internal/github"
)

func TestEvaluate_Empty(t *testing.T) {
	l := Evaluate(&github.Stats{})
	if l.Any() {
		t.Fatalf("expected empty loadout, got %+v", l)
	}
}

func TestEvaluate_RanksByCommitCount(t *testing.T) {
	l := Evaluate(&github.Stats{Languages: map[string]int{
		"Python": 10,
		"Go":     100,
		"Rust":   50,
	}})
	if l.Weapon == nil || l.Weapon.Lang != "Go" {
		t.Fatalf("expected Go weapon, got %+v", l.Weapon)
	}
	if l.Shield == nil || l.Shield.Lang != "Rust" {
		t.Fatalf("expected Rust shield, got %+v", l.Shield)
	}
	if l.Accessory == nil || l.Accessory.Lang != "Python" {
		t.Fatalf("expected Python accessory, got %+v", l.Accessory)
	}
}

func TestEvaluate_SkipsUnknownLanguages(t *testing.T) {
	l := Evaluate(&github.Stats{Languages: map[string]int{
		"COBOL": 1000, // not in items map
		"Go":    10,
	}})
	if l.Weapon == nil || l.Weapon.Lang != "Go" {
		t.Fatalf("expected Go weapon (unknown lang skipped), got %+v", l.Weapon)
	}
	if l.Shield != nil {
		t.Fatalf("expected no shield, got %+v", l.Shield)
	}
}

func TestEvaluate_MoreThanThreeLanguagesTruncated(t *testing.T) {
	l := Evaluate(&github.Stats{Languages: map[string]int{
		"Go":         40,
		"Rust":       30,
		"Python":     20,
		"JavaScript": 10,
	}})
	if l.Weapon.Lang != "Go" || l.Shield.Lang != "Rust" || l.Accessory.Lang != "Python" {
		t.Fatalf("unexpected loadout: %+v %+v %+v", l.Weapon, l.Shield, l.Accessory)
	}
}

func TestRankLanguages_TiesBreakAlphabetically(t *testing.T) {
	got := rankLanguages(map[string]int{"Rust": 5, "Go": 5, "Python": 5})
	want := []string{"Go", "Python", "Rust"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestLoadout_Slots(t *testing.T) {
	l := Evaluate(&github.Stats{Languages: map[string]int{"Go": 1}})
	slots := l.Slots()
	if len(slots) != 3 {
		t.Fatalf("expected 3 slots, got %d", len(slots))
	}
	if slots[0].Label != "Weapon" || slots[0].Item == nil || slots[0].Item.Lang != "Go" {
		t.Fatalf("unexpected weapon slot: %+v", slots[0])
	}
	if slots[1].Item != nil || slots[2].Item != nil {
		t.Fatalf("expected empty shield/accessory, got %+v %+v", slots[1], slots[2])
	}
}
