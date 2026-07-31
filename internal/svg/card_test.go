package svg

import (
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/stats"
)

func TestCard_EscapesLogin(t *testing.T) {
	malicious := `<script>alert(1)</script>`
	char := &stats.Character{Class: "Guardian", Title: "The Adventurer"}
	svg, err := Card(malicious, char, "")
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
