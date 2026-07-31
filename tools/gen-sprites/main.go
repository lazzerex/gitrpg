// gen-sprites renders the 16x16 class character sprites from text grids into
// web/static/assets/sprites/. Run: go run ./tools/gen-sprites
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func hex(s string) color.NRGBA {
	var r, g, b uint8
	_, err := fmtSscanf(s, &r, &g, &b)
	if err != nil {
		log.Fatalf("bad hex %q: %v", s, err)
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func fmtSscanf(s string, r, g, b *uint8) (int, error) {
	var ri, gi, bi int
	n, err := sscanfHex(s, &ri, &gi, &bi)
	*r, *g, *b = uint8(ri), uint8(gi), uint8(bi)
	return n, err
}

func sscanfHex(s string, r, g, b *int) (int, error) {
	s = strings.TrimPrefix(s, "#")
	v, err := parseHex(s)
	if err != nil {
		return 0, err
	}
	*r = v >> 16 & 0xff
	*g = v >> 8 & 0xff
	*b = v & 0xff
	return 3, nil
}

func parseHex(s string) (int, error) {
	v := 0
	for _, c := range s {
		v <<= 4
		switch {
		case c >= '0' && c <= '9':
			v |= int(c - '0')
		case c >= 'a' && c <= 'f':
			v |= int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			v |= int(c-'A') + 10
		default:
			return 0, os.ErrInvalid
		}
	}
	return v, nil
}

// darken returns the color at 55% brightness for shading.
func darken(c color.NRGBA) color.NRGBA {
	return color.NRGBA{R: uint8(int(c.R) * 55 / 100), G: uint8(int(c.G) * 55 / 100), B: uint8(int(c.B) * 55 / 100), A: 255}
}

var classAccents = map[string]string{
	"guardian":   "#00add8",
	"berserker":  "#e05d44",
	"paladin":    "#3178c6",
	"rogue":      "#e8c94a",
	"sage":       "#4b8bbe",
	"knight":     "#9b72cf",
	"battlemage": "#c07d28",
	"warlord":    "#f34b7d",
	"wanderer":   "#6e7681",
}

// Shared body, rows 8-15. Runes: . transparent, O outline, S skin,
// A armor (accent), D armor shade, G gold trim, L leg, B boot.
var body = []string{
	"...OAAAAAAAAO...",
	"..OSAADAADAASO..",
	"..OSADAAAADASO..",
	"..OSAAAAAAAASO..",
	"...OGGGGGGGGO...",
	"....OLL..LLO....",
	"....OLL..LLO....",
	"....OBB..BBO....",
}

// bodyStride is the second walk frame: legs apart.
var bodyStride = []string{
	"...OAAAAAAAAO...",
	"..OSAADAADAASO..",
	"..OSADAAAADASO..",
	"..OSAAAAAAAASO..",
	"...OGGGGGGGGO...",
	"...OLL....LLO...",
	"...OLL....LLO...",
	"...OBB....BBO...",
}

// Per-class heads, rows 0-7. Extra runes: E eye, H headgear (accent),
// D headgear shade, G gold, X hair (neutral dark).
var heads = map[string][]string{
	// full helm with open visor
	"guardian": {
		"................",
		".....OOOOOO.....",
		"....OHHHHHHO....",
		"...OHHHHHHHHO...",
		"...OHDSSSSDHO...",
		"...OHSESSESHO...",
		"...OHSSSSSSHO...",
		"....OSSSSSSO....",
	},
	// horned helm
	"berserker": {
		"..O..........O..",
		".OHO..OOOO..OHO.",
		".OHHOOHHHHOOHHO.",
		"..OHHHHHHHHHHO..",
		"...OHSSSSSSHO...",
		"...OSSESSESSO...",
		"...OSSSSSSSSO...",
		"....OSSSSSSO....",
	},
	// circlet with gem
	"paladin": {
		"................",
		".....OOOOOO.....",
		"....OXXXXXXO....",
		"...OGGGDGGGGO...",
		"...OXSSSSSSXO...",
		"...OSSESSESSO...",
		"...OSSSSSSSSO...",
		"....OSSSSSSO....",
	},
	// hood, shadowed face
	"rogue": {
		"................",
		".....OOOOOO.....",
		"....OHHHHHHO....",
		"...OHHHHHHHHO...",
		"...OHDDDDDDHO...",
		"...OHDESSEDHO...",
		"...OHDSSSSDHO...",
		"....ODSSSSDO....",
	},
	// wide-brim wizard hat
	"sage": {
		".......OO.......",
		"......OHHO......",
		".....OHHHHO.....",
		"....OHHHHHHO....",
		".OOOHHHHHHHHOOO.",
		"OHHHHHHHHHHHHHHO",
		"..OOSSESSESOO...",
		"....OSSSSSSO....",
	},
	// plumed helm
	"knight": {
		".......OGO......",
		".....OOOGOO.....",
		"....OHHHGHHO....",
		"...OHHHHHHHHO...",
		"...OHDSSSSDHO...",
		"...OHSESSESHO...",
		"...OHSSSSSSHO...",
		"....OSSSSSSO....",
	},
	// pointed hat with gold band
	"battlemage": {
		".......OO.......",
		"......OHHO......",
		".....OHHHHO.....",
		"....OGGGGGGO....",
		".OOOHHHHHHHHOOO.",
		"OHHHHHHHHHHHHHHO",
		"..OOSSESSESOO...",
		"....OSSSSSSO....",
	},
	// spiked crown helm
	"warlord": {
		"...O...OO...O...",
		"..OGO.OGGO.OGO..",
		"..OGOOHHHHOOGO..",
		"...OHHHHHHHHO...",
		"...OHDSSSSDHO...",
		"...OHSESSESHO...",
		"...OHSSSSSSHO...",
		"....OSSSSSSO....",
	},
	// headband over hair
	"wanderer": {
		"................",
		".....OOOOOO.....",
		"....OXXXXXXO....",
		"...OHHHHHHHHO...",
		"...OXSSSSSSXO...",
		"...OSSESSESSO...",
		"...OSSSSSSSSO...",
		"....OSSSSSSO....",
	},
}

// slime is the minigame enemy. Own palette: G body, D shade, W eye white, B pupil.
var slime = []string{
	"................",
	"................",
	"................",
	"................",
	"......OOOO......",
	"....OOGGGGOO....",
	"...OGGGGGGGGO...",
	"..OGGWGGGGWGGO..",
	"..OGGBGGGGBGGO..",
	".OGGGGGGGGGGGGO.",
	".OGGGGGGGGGGGGO.",
	"OGGGGGGGGGGGGGGO",
	"OGGDGGGDDGGGDGGO",
	"OGDDGGDDDDGGDDGO",
	".OODDDDDDDDDDOO.",
	"..OOOOOOOOOOOO..",
}

var slimePalette = map[rune]color.NRGBA{
	'O': {R: 0x1a, G: 0x12, B: 0x26, A: 255},
	'G': {R: 0x5c, G: 0xb8, B: 0x5c, A: 255},
	'D': {R: 0x2d, G: 0x6a, B: 0x2d, A: 255},
	'W': {R: 0xff, G: 0xff, B: 0xff, A: 255},
	'B': {R: 0x1a, G: 0x12, B: 0x26, A: 255},
}

// imp is the ranged enemy: red, horned.
var imp = []string{
	"................",
	"..O.........O...",
	".OHO.......OHO..",
	".OHHO.OOO.OHHO..",
	"..OHHOHHHOHHO...",
	"...OHHHHHHHO....",
	"..OHWHHHHWHHO...",
	"..OHBHHHHBHHO...",
	"..OHHHHHHHHHO...",
	"...OHHDDDHHO....",
	"...OHHHHHHHO....",
	"..OHHOHHHOHHO...",
	"..OHO.OHO.OHO...",
	"...O.OHHHO.O....",
	".....OHOHO......",
	"......O.O.......",
}

// wraith is the teleporting enemy: purple ghost.
var wraith = []string{
	"................",
	".....OOOOOO.....",
	"....OHHHHHHO....",
	"...OHHHHHHHHO...",
	"..OHWHHHHWHHHO..",
	"..OHBHHHHBHHHO..",
	"..OHHHHHHHHHHO..",
	"..OHHDDDDHHHHO..",
	"..OHHHHHHHHHHO..",
	"..OHHHHHHHHHHO..",
	"..OHHHHHHHHHHO..",
	"..OHHOHHOHHOHO..",
	"..OHO.OHO.OHO...",
	"..OO...O...OO...",
	"................",
	"................",
}

// Pickups: R heart red, G gold, C cyan, W white highlight.
var heart = []string{
	"................",
	"................",
	"...OO....OO.....",
	"..ORRO..ORRO....",
	".ORRRROORRRRO...",
	".ORWRRRRRRRRO...",
	".ORWRRRRRRRRO...",
	".ORRRRRRRRRRO...",
	"..ORRRRRRRRO....",
	"...ORRRRRRO.....",
	"....ORRRRO......",
	".....ORRO.......",
	"......OO........",
	"................",
	"................",
	"................",
}

var star = []string{
	"................",
	".......OO.......",
	"......OGGO......",
	"......OGGO......",
	".....OGGGGO.....",
	".OOOOOGGGGOOOOO.",
	".OGGGGGWGGGGGGO.",
	"..OGGGGGGGGGGO..",
	"...OGGGGGGGGO...",
	"....OGGGGGGO....",
	"....OGGGGGGO....",
	"...OGGOOOOGGO...",
	"...OGO....OGO...",
	"...OO......OO...",
	"................",
	"................",
}

var bolt = []string{
	"................",
	"......OOOO......",
	".....OCCCCO.....",
	"....OCCCCO......",
	"...OCCCCO.......",
	"...OCCCCOOOO....",
	"..OCCCCCCCCO....",
	"..OOOOCCCCO.....",
	".....OCCCO......",
	"....OCCCO.......",
	"...OCCCO........",
	"...OCCO.........",
	"..OCCO..........",
	"..OCO...........",
	"..OO............",
	"................",
}

// boss is the final-wave Warden: crowned demon, drawn at 64px in-game.
var boss = []string{
	"..O..........O..",
	".OHO..OOOO..OHO.",
	".OHHOOGGGGOOHHO.",
	"..OHHGGGGGGHHO..",
	"..OHHHHHHHHHHO..",
	".OHWWHHHHWWHHO..",
	".OHBWHHHHBWHHO..",
	".OHHHHHHHHHHHO..",
	".OHHDDDDDDDHHO..",
	".OHHDHHHHHDHHO..",
	".OHHHHHHHHHHHO..",
	".OHHOHHHHHOHHO..",
	".OHO.OHHHO.OHO..",
	"..O..OHHHO..O...",
	".....OHOHO......",
	"......O.O.......",
}

var outlineC = color.NRGBA{R: 0x1a, G: 0x12, B: 0x26, A: 255}
var whiteC = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255}

var impPalette = map[rune]color.NRGBA{
	'O': outlineC, 'W': whiteC, 'B': outlineC,
	'H': {R: 0xe0, G: 0x5d, B: 0x44, A: 255},
	'D': {R: 0x7b, G: 0x33, B: 0x25, A: 255},
}

var wraithPalette = map[rune]color.NRGBA{
	'O': outlineC, 'W': whiteC, 'B': outlineC,
	'H': {R: 0x9b, G: 0x72, B: 0xcf, A: 255},
	'D': {R: 0x55, G: 0x3e, B: 0x72, A: 255},
}

var bossPalette = map[rune]color.NRGBA{
	'O': outlineC, 'W': whiteC, 'B': outlineC,
	'H': {R: 0x8e, G: 0x2a, B: 0x3c, A: 255},
	'D': {R: 0x4d, G: 0x16, B: 0x21, A: 255},
	'G': {R: 0xff, G: 0xd7, B: 0x00, A: 255},
}

var pickupPalette = map[rune]color.NRGBA{
	'O': outlineC, 'W': whiteC,
	'R': {R: 0xe0, G: 0x5d, B: 0x44, A: 255},
	'G': {R: 0xff, G: 0xd7, B: 0x00, A: 255},
	'C': {R: 0x00, G: 0xad, B: 0xd8, A: 255},
}

func palette(accentHex string) map[rune]color.NRGBA {
	accent := hex(accentHex)
	return map[rune]color.NRGBA{
		'O': hex("#1a1226"),
		'S': hex("#e8b88a"),
		'E': hex("#1a1226"),
		'A': accent,
		'H': accent,
		'D': darken(accent),
		'G': hex("#ffd700"),
		'X': hex("#4a3626"),
		'L': darken(accent),
		'B': hex("#3d2b1f"),
	}
}

func renderSprite(head []string, bodyRows []string, pal map[rune]color.NRGBA) *image.NRGBA {
	rows := append(append([]string{}, head...), bodyRows...)
	return renderGrid(rows, pal)
}

func renderGrid(rows []string, pal map[rune]color.NRGBA) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	for y, row := range rows {
		for x, r := range row {
			if r == '.' {
				continue
			}
			c, ok := pal[r]
			if !ok {
				log.Fatalf("unknown rune %q at %d,%d", r, x, y)
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

func scale(src *image.NRGBA, factor int) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx()*factor, b.Dy()*factor))
	for y := 0; y < dst.Bounds().Dy(); y++ {
		for x := 0; x < dst.Bounds().Dx(); x++ {
			dst.SetNRGBA(x, y, src.NRGBAAt(x/factor, y/factor))
		}
	}
	return dst
}

func writePNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func main() {
	outDir := "web/static/assets/sprites"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	order := []string{"guardian", "berserker", "paladin", "rogue", "sage", "knight", "battlemage", "warlord", "wanderer"}

	const previewScale = 8
	sheet := image.NewNRGBA(image.Rect(0, 0, len(order)*(16*previewScale+8)+8, 16*previewScale+16))
	draw.Draw(sheet, sheet.Bounds(), &image.Uniform{hex("#050010")}, image.Point{}, draw.Src)

	for i, class := range order {
		head, ok := heads[class]
		if !ok {
			log.Fatalf("no head for %s", class)
		}
		pal := palette(classAccents[class])
		sprite := renderSprite(head, body, pal)
		writePNG(filepath.Join(outDir, class+".png"), sprite)
		writePNG(filepath.Join(outDir, class+"_b.png"), renderSprite(head, bodyStride, pal))

		big := scale(sprite, previewScale)
		off := image.Pt(8+i*(16*previewScale+8), 8)
		draw.Draw(sheet, big.Bounds().Add(off), big, image.Point{}, draw.Over)
	}

	writePNG(filepath.Join(outDir, "slime.png"), renderGrid(slime, slimePalette))
	writePNG(filepath.Join(outDir, "imp.png"), renderGrid(imp, impPalette))
	writePNG(filepath.Join(outDir, "wraith.png"), renderGrid(wraith, wraithPalette))
	writePNG(filepath.Join(outDir, "heart.png"), renderGrid(heart, pickupPalette))
	writePNG(filepath.Join(outDir, "star.png"), renderGrid(star, pickupPalette))
	writePNG(filepath.Join(outDir, "bolt.png"), renderGrid(bolt, pickupPalette))
	writePNG(filepath.Join(outDir, "boss.png"), renderGrid(boss, bossPalette))

	writePNG(filepath.Join(outDir, "preview.png"), sheet)
	log.Printf("wrote %d class sprites + 3 enemies + 3 pickups + preview to %s", len(order), outDir)
}
