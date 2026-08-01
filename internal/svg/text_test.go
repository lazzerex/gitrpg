package svg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/equipment"
	"github.com/lazzerex/gitrpg/internal/stats"
)

// withGlyphs loads the real card font and restores the unloaded state
// afterwards, so the fallback-renderer tests stay unaffected.
func withGlyphs(t *testing.T) {
	t.Helper()
	if err := loadGlyphs("../../web/static/assets/fonts/PressStart2P-Regular.ttf"); err != nil {
		t.Fatalf("loadGlyphs: %v", err)
	}
	t.Cleanup(func() { glyphs = nil })
}

func TestLoadGlyphs_CoversPrintableASCII(t *testing.T) {
	withGlyphs(t)
	for r := rune(' '); r <= '~'; r++ {
		if _, ok := glyphs[r]; !ok {
			t.Errorf("missing glyph for %q", r)
		}
	}
	if got := glyphs['A'].advance; got != 1 {
		t.Errorf("advance = %v, want 1 em (monospace)", got)
	}
	if len(glyphs['A'].cmds) == 0 {
		t.Error("glyph A has no outline commands")
	}
}

func TestTextWidth_ScalesWithSize(t *testing.T) {
	withGlyphs(t)
	if got := textWidth("ABC", 8); got != 24 {
		t.Errorf("textWidth = %v, want 24", got)
	}
}

func TestTextPath_Anchors(t *testing.T) {
	withGlyphs(t)
	startX := func(anchor string) float64 {
		var b strings.Builder
		textPath(&b, 100, 50, 8, "#fff", anchor, "ABC")
		var x float64
		if _, err := fmt.Sscanf(b.String(), `<path d="M%g`, &x); err != nil {
			t.Fatalf("parse path start: %v", err)
		}
		return x
	}
	left, mid, end := startX(""), startX("middle"), startX("end")
	if mid != left-12 {
		t.Errorf("middle anchor starts at %v, want %v", mid, left-12)
	}
	if end != left-24 {
		t.Errorf("end anchor starts at %v, want %v", end, left-24)
	}
}

func TestCard_RendersTextAsPaths(t *testing.T) {
	withGlyphs(t)
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer", Level: 27, XPForLevel: 100}
	svg, err := Card("tester", char, equipment.Loadout{})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if strings.Contains(svg, "<text") {
		t.Error("card still emits <text>; webfont-dependent rendering")
	}
	if strings.Contains(svg, "@font-face") {
		t.Error("card still embeds @font-face")
	}
	if !strings.Contains(svg, `<path d="M`) {
		t.Error("card has no glyph paths")
	}
}

func TestCard_LoginRendersAsGlyphsNotMarkup(t *testing.T) {
	withGlyphs(t)
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer", XPForLevel: 100}
	svg, err := Card(`<script>alert(1)</script>`, char, equipment.Loadout{})
	if err != nil {
		t.Fatalf("Card() error: %v", err)
	}
	if strings.Contains(svg, "<script>") {
		t.Errorf("card output contains unescaped <script> tag:\n%s", svg)
	}
}
