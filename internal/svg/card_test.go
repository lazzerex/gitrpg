package svg

import (
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/equipment"
	"github.com/lazzerex/gitrpg/internal/stats"
)

func TestCard_EscapesLogin(t *testing.T) {
	malicious := `<script>alert(1)</script>`
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer"}
	svg, err := Card(malicious, char, equipment.Loadout{})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if strings.Contains(svg, "<script>") {
		t.Errorf("Card() output contains unescaped <script> tag:\n%s", svg)
	}
	if !strings.Contains(svg, "&lt;script&gt;") {
		t.Errorf("Card() output missing escaped login, want &lt;script&gt;:\n%s", svg)
	}
}

func TestDemo_EscapesClassParam(t *testing.T) {
	malicious := `<script>alert(1)</script>`
	svg, err := Demo(malicious, 0)
	if err != nil {
		t.Fatalf("Demo() error: %v", err)
	}
	if strings.Contains(svg, "<script>") {
		t.Errorf("Demo() output contains unescaped <script> tag:\n%s", svg)
	}
}

func TestCard_GearSlots(t *testing.T) {
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer", XPForLevel: 100}

	empty, err := Card("tester", char, equipment.Loadout{})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if got := strings.Count(empty, ">EMPTY</text>"); got != 3 {
		t.Errorf("empty loadout: %d EMPTY labels, want 3", got)
	}

	weapon, _ := equipment.Lookup("go-compiler")
	svg, err := Card("tester", char, equipment.Loadout{Weapon: &weapon})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if !strings.Contains(svg, "GO COMPILER") {
		t.Error("loadout with weapon missing item name")
	}
	if got := strings.Count(svg, ">EMPTY</text>"); got != 2 {
		t.Errorf("one-item loadout: %d EMPTY labels, want 2", got)
	}
}

func TestCard_MaxStatGold(t *testing.T) {
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer", Intelligence: 100, XPForLevel: 100}
	svg, err := Card("tester", char, equipment.Loadout{})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if !strings.Contains(svg, ">MAX</text>") {
		t.Error("maxed stat missing MAX label")
	}
}
