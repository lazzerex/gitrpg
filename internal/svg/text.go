package svg

import (
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Card text is drawn as glyph outlines rather than <text> with an embedded
// @font-face: an SVG loaded as an image (README embed, GitHub camo) renders in
// secure static mode, where webfonts are unreliable and the card silently falls
// back to a system font.

type pathCmd struct {
	op  byte
	pts [3][2]float64
}

// glyph holds outline commands and advance width in em units, y-down from the baseline.
type glyph struct {
	cmds    []pathCmd
	advance float64
}

var glyphs map[rune]glyph

// loadGlyphs parses a TrueType file and caches printable ASCII outlines.
func loadGlyphs(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	f, err := sfnt.Parse(data)
	if err != nil {
		return err
	}

	upem := float64(f.UnitsPerEm())
	ppem := fixed.I(int(f.UnitsPerEm()))
	unit := func(v fixed.Int26_6) float64 { return float64(v) / 64 / upem }

	var buf sfnt.Buffer
	out := make(map[rune]glyph, 96)
	for r := rune(' '); r <= '~'; r++ {
		idx, err := f.GlyphIndex(&buf, r)
		if err != nil || idx == 0 {
			continue
		}
		adv, err := f.GlyphAdvance(&buf, idx, ppem, font.HintingNone)
		if err != nil {
			continue
		}
		segs, err := f.LoadGlyph(&buf, idx, ppem, nil)
		if err != nil {
			continue
		}
		g := glyph{advance: unit(adv)}
		for _, s := range segs {
			c := pathCmd{}
			switch s.Op {
			case sfnt.SegmentOpMoveTo:
				c.op = 'M'
			case sfnt.SegmentOpLineTo:
				c.op = 'L'
			case sfnt.SegmentOpQuadTo:
				c.op = 'Q'
			case sfnt.SegmentOpCubeTo:
				c.op = 'C'
			}
			for i, p := range s.Args {
				c.pts[i] = [2]float64{unit(p.X), unit(p.Y)}
			}
			g.cmds = append(g.cmds, c)
		}
		out[r] = g
	}
	if len(out) == 0 {
		return os.ErrInvalid
	}
	glyphs = out
	return nil
}

func textWidth(s string, size float64) float64 {
	var w float64
	for _, r := range s {
		if g, ok := glyphs[r]; ok {
			w += g.advance * size
		}
	}
	return w
}

func num(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// glyphPath appends one glyph's outline, positioned at baseline point (x, y).
func glyphPath(b *strings.Builder, g glyph, x, y, size float64) {
	open := false
	for _, c := range g.cmds {
		if c.op == 'M' {
			if open {
				b.WriteByte('Z')
			}
			open = true
		}
		b.WriteByte(c.op)
		n := 1
		switch c.op {
		case 'Q':
			n = 2
		case 'C':
			n = 3
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(num(x + c.pts[i][0]*size))
			b.WriteByte(' ')
			b.WriteString(num(y + c.pts[i][1]*size))
		}
	}
	if open {
		b.WriteByte('Z')
	}
}

// textPath renders s as a single filled path. anchor is "", "middle" or "end",
// matching the text-anchor attribute it replaces.
func textPath(b *strings.Builder, x, y, size float64, fill, anchor, s string) {
	switch anchor {
	case "middle":
		x -= textWidth(s, size) / 2
	case "end":
		x -= textWidth(s, size)
	}

	var d strings.Builder
	for _, r := range s {
		g, ok := glyphs[r]
		if !ok {
			continue
		}
		glyphPath(&d, g, x, y, size)
		x += g.advance * size
	}
	if d.Len() == 0 {
		return
	}
	b.WriteString(`<path d="`)
	b.WriteString(d.String())
	b.WriteString(`" fill="`)
	b.WriteString(fill)
	b.WriteString(`"/>`)
}
