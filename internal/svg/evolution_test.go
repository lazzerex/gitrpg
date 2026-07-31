package svg

import (
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/stats"
)

func TestCardStage(t *testing.T) {
	tests := []struct{ level, want int }{
		{1, 0}, {9, 0}, {10, 1}, {24, 1}, {25, 2}, {49, 2}, {50, 3}, {99, 3},
	}
	for _, tt := range tests {
		if got := cardStage(tt.level); got != tt.want {
			t.Errorf("cardStage(%d) = %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestCard_EvolutionTiers(t *testing.T) {
	render := func(level int) string {
		char := &stats.Character{Class: "Guardian", Level: level, XPForLevel: 100}
		out, err := Card("tester", char, "")
		if err != nil {
			t.Fatalf("Card(level %d): %v", level, err)
		}
		return out
	}

	base := render(5)
	for _, marker := range []string{"VETERAN", "ELITE", "LEGENDARY", "<animate"} {
		if strings.Contains(base, marker) {
			t.Errorf("level 5 card unexpectedly contains %q", marker)
		}
	}
	if !strings.Contains(base, tierBorders[0]) {
		t.Error("level 5 card missing bronze frame")
	}

	veteran := render(10)
	if !strings.Contains(veteran, "VETERAN") {
		t.Error("level 10 card missing VETERAN label")
	}
	if !strings.Contains(veteran, tierBorders[1]) {
		t.Error("level 10 card missing silver frame")
	}
	if strings.Contains(veteran, "ELITE") || strings.Contains(veteran, "<animate") {
		t.Error("level 10 card has tiers above veteran")
	}

	elite := render(25)
	if !strings.Contains(elite, "ELITE") {
		t.Error("level 25 card missing ELITE label")
	}
	if !strings.Contains(elite, tierBorders[2]) {
		t.Error("level 25 card missing gold frame")
	}
	if strings.Contains(elite, "<animate") {
		t.Error("level 25 card has legendary animation")
	}

	legendary := render(50)
	if !strings.Contains(legendary, "LEGENDARY") {
		t.Error("level 50 card missing LEGENDARY label")
	}
	if !strings.Contains(legendary, `stroke-width="4"`) {
		t.Error("level 50 card missing thick gold frame")
	}
	if !strings.Contains(legendary, "<animate attributeName=\"stroke-opacity\"") ||
		!strings.Contains(legendary, "animateTransform") {
		t.Error("level 50 card missing shimmer animations")
	}
}
