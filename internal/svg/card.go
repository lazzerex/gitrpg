package svg

import (
	"encoding/base64"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazzerex/gitrpg/internal/equipment"
	"github.com/lazzerex/gitrpg/internal/stats"
)

func ClassColor(class string) string {
	if c, ok := classColors[class]; ok {
		return c
	}
	return classColors["Wanderer"]
}

var classColors = map[string]string{
	"Berserker":  "#e05d44",
	"Guardian":   "#00add8",
	"Paladin":    "#3178c6",
	"Rogue":      "#e8c94a",
	"Sage":       "#4b8bbe",
	"Knight":     "#9b72cf",
	"Battlemage": "#c07d28",
	"Warlord":    "#f34b7d",
	"Wanderer":   "#6e7681",
}

// Card palette: muted night slate with a single class accent.
const (
	colBG      = "#0b0d1a"
	colWall    = "#181a2e"
	colWallLt  = "#20233c"
	colFloorA  = "#232642"
	colFloorB  = "#1a1c30"
	colBone    = "#dcd7c9"
	colFrameD  = "#565d78"
	colText    = "#e9e4d0"
	colMuted   = "#7f8aa3"
	colGold    = "#e3b458"
	colGoldD   = "#8a6a30"
	colSilver  = "#c0c8d0"
	colBarBG   = "#1e2136"
	colFlameO  = "#e8722c"
	colFlameY  = "#f5c542"
	colShadow  = "#0a0b16"
	colDivider = "#20243c"
)

const (
	cardW = 700
	cardH = 308
)

// 32×32 pixel art item icons (base64 PNG).
const (
	iconSword  = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABLklEQVRYR9WWMQ6CMBSGYdRbOABx4gxM3oTdWTa5gd6EhOhE4uLMiCwOnMJF428e4ZFSQCmtLH/6StLv/f1bsC3Lelrix+6oT1p+L2IGQFEUrDPP82is1InaAe0ASZIwB1zXxVi1E7UDOgGoc4QxjmOMfd+HqnaiGTDtAFqcEB2xWZ2QnfFZQIwGmCUTQ65ZpVsxBOBbJ9ofOeFafwUgdSLPc8xHUQQNwxCaZRm0LEvh13WMA8YACEEWmx3qx9MNur5uWedTOmAMAANZ7T9/VI/LAbq8n6V7LwwEFUcqjpsOACxMaa+qqt0x60NZBrQDOI4j7JTqQRBgnhxK05Rt/zf3AAufTgAG0rCBmkJGCJCUskD6iwPGAPSdWuYEvTylA8YDdG0V6lNkoM8BKcAL+MPCG+UmDIIAAAAASUVORK5CYIIA"
	iconSword2 = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABUklEQVRYR82WO26DQBCGQ5sz0LgAlIpbuOAguE5PZ05BR02PlN5SmjSULFKUgjqda1v5pVE8y7I2j/VAM+yygm++2Qfey/91ubn/u/W0tpPm7Ue2AdC2LTKNoogydmpiYGAzAJS+axOjBsQBlFJgCMPQ6ZwYNbAZACqFKxN3DYgBVFWFb8dxzHbAtU2MGpAAoEyxJdOG5BrEtM2KA6xlQj/cjGeL7aBZamIxwCQTTdNgfJZliGmastVTFMVkA5sBeAiEDJRlifFBECB2Xcei/qc15WfDOieeAcBM5HnOavy1O6DdfP8ivn2+I9Z1bay9tZO9ediACUkAZoJq7e2P6D+rE+Lrz4e19ksMiAOgBLTO+75nRaJZT51kaGwuTFkFLHNxAMqMqPTMkyTBI9/3EcmUbmK2AXEAw1LVk8FcsYBi/GwDkgB39qnBY+uxPMfAqgBX+oXuIQYvYaIAAAAASUVORK5CYIIA"
	iconGem    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABRklEQVRYR9WXMQrCUAxAFXXp4iAIXXX2AN5AwSN08wJOHsJDuLSewd3FzcXZVa1bh0qhiELsB/tNfxIHoy6B337z8pL+0npN+VdXzl/7W4B7hTlxQeINRWI1AEh87PWAwz8cIC59H+L0fDZi2IWxb3ytXAOgVHl3OASeRhRBPPX7ENfXq9gE14AaAJrYNNo2YNYlJigDagDOxJQBiYkqA2oAvMSj0bPIIIBwK+Jlu0UPSNdM2AZ+A+AehryXpGWA2tRcrd5OStSAOsCs0wHSQauFFjX2vNK66fE+z9H7N1kG67s05RnQBDCEMIxckEWSiCt/U2H9gzoAy4SZBduAq+e2Kta7oKoV3wBwmrABJJVTM2CbQmfimwCoiXm7DeuTOK58zqnTkZoBpwkNgJIJpDppQR9/Gal9F1AtFV8XKxNnIDaoAzwAtSfaIXzqsLAAAAAASUVORK5CYIIA"
	iconKey    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAA7klEQVRYR+2WwQ3DIAxFyTIwDWuwDQuwBtPAMq2E5IMtN8VAMVKTC5IFyfPnx/Zlxp/Xh6OX5JWizeTFagDtw6UUNlHnHMS7kuvaxGWuAYAyr7U2Lu99W3POiBPixpjbJCUKnAXA3HUDBCV+roAaAGQIa0oJeYDGl3tAEwBcju7aWtviMca2giLfMhcVC64OgBIaAEiJEMJQ5jMKnAUg/e9pA5FUQnp2qPA8ACsUQIOIhgfOAICBhM4FvRVwpg7czgU7ANie0Nv/V5hQDYAdw3f2gjMAqPvhPnZ4gHX/dgDq4r/yANuOaXB7IZoFeAN+/6whBbi1kAAAAABJRU5ErkJggg=="
	iconPotion = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABIElEQVRYR2NkIB78J14pWCUjMeqJUgQ1aHA4wEdPEuwef0tNFA9uPH4dzN9y6TlMnCjPEaUIOQQGwgEoQU+CA4gKCWJCYHA4YOVWiDu4/m0B0z4+PihpYMsWiPg3Joh4uDfcb3g9SXQIDBoHHD75CuxDW3MxlBBAF6dZCAyYA2DehUVFmBdqObdqG4SP5HPa5IKBcADMJ+BsQIIDiEngxFUYyCXhoHHAx6VXUBIBf7QOSfmfqASCmswYUKJgwB1wtfYwivu0m23pGwKjDhg0IXCTcTs47oOa2gYmDQyYA9bVVaHkArqHwKgDBl0IiHoXgNPElCZxkop5oqpMbLUhegjQzQE5dS/R6ihULs1DYCAdgNIywhMMpEQrSS0imjgAAPoJziGHZTDaAAAAAElFTkSuQmCC"
	iconShield = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABmElEQVRYR2NkGGDASKL9/4lUT7S5RCuEWkx3B6BY2Hsau/2Hb71BCZgN0aLoAYXTo4RCYMAcALYY5mOYD03uzwT77NGjR1iTgpycHFj8jGI6rhDB8DCuEBgwB4AtDlj6GuwDWzURMD0jSg1Mq6qqovhMRkYGzH/y5AmK+O3bt8F87YZjBEMCPQQGlwOuNlih+IDUEIBphoUEUu6AexxvCNDTAeCgV8teAXZ0nOQdMH3sGGocwuI2Pj4eJTfAUv/ChQuxphUrK0hILnquAqZvTY2ABQ4jLAQGlwMcfu/DmrphPoH5FBa3sKiChQx6yMFyywFWJ+JCYNA4ACULIHFe2beilBMzF6wE82EOx6WP5BAYcAegl3DoaYDRrRnsxv+7asE03dMAzRwA8xF6yUe3cmAgHABLbygFEswhMEly6wJYVCGXgPCiEC2lDy4HwIMFmspJDQGYz2HmkBwCA+EArGkBJgirJQm1CWG1Hj6f40oDg8YBKA6BcWDtBVxFNBYfE/IoA0n9goFwALpn6d41G3QOIBT9JMsDANCtYDCjzkYBAAAAAElFTkSuQmCC"
	iconStaff  = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAiElEQVR4nGNgGOmAEZfErKUb/1PDgrRof5x2MDAwMDBRwxJKwKgDWAgpOHTwCFkG29nbEKVu8IcAsT4hFwx4CODNo/hAV5QxSjlRtuwsWWYNeAiMOmDUAQTLAVzgybuvVHHA0A0BGSFuqjhgwENg1AGjuWA0F5Bch6O3A9ABqe2CAQ+BUQeMAgA30RX7GlRZ2QAAAABJRU5ErkJggg=="
	iconHelmet = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAtklEQVR4nGNgGAUEwKylG//PWrrxP7F8UgETpQ6kFDDikiDXV2nR/jjNxAYGPARYaG0Bekiih9DgD4FDB48wMDAwMNjZ2xDFhwFi09CAhwDVcwEugCt3DHgI0DwXECoXhm8IEFsiDngIEHQlqblhyNUFow7AmQseb4n4z8DAwLD9I2kGwvTJ+qwYormAXnUADAx4COBMA7B6nlyA3j7ABQZvCBDrA0rB4AsBT/7ldHXAgIfAKAAAwcFBPKy3LrwAAAAASUVORK5CYII="
)

// gearIcons maps equipment item slugs to pixel icons.
var gearIcons = map[string]string{
	"go-compiler":       iconSword,
	"rust-axe":          iconSword2,
	"typescript-lens":   iconGem,
	"javascript-dagger": iconKey,
	"python-flask":      iconPotion,
	"csharp-bulwark":    iconShield,
	"java-hammer":       iconStaff,
	"cpp-gauntlet":      iconHelmet,
}

var spriteB64 = map[string]string{}

// LoadAssets reads class sprites and the card font from the static assets
// directory. Call once at startup, before serving cards.
func LoadAssets(assetsDir string) error {
	var firstErr error
	for class := range classColors {
		path := filepath.Join(assetsDir, "sprites", strings.ToLower(class)+".png")
		data, err := os.ReadFile(path)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		spriteB64[class] = base64.StdEncoding.EncodeToString(data)
	}
	if err := loadGlyphs(filepath.Join(assetsDir, "fonts", "PressStart2P-Regular.ttf")); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

// cardStage maps level to an evolution tier per GAME_DESIGN.md:
// visual card upgrades at levels 10, 25 and 50.
func cardStage(level int) int {
	switch {
	case level >= 50:
		return 3
	case level >= 25:
		return 2
	case level >= 10:
		return 1
	default:
		return 0
	}
}

var tierLabels = [4]string{"", "VETERAN", "ELITE", "LEGENDARY"}

// tierColors is the tier banner and scene border color per stage.
var tierColors = [4]string{colFrameD, colSilver, colGold, colGold}

func rct(b *strings.Builder, x, y, w, h int, fill string) {
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x, y, w, h, fill)
}

func txt(b *strings.Builder, x, y, size int, fill, anchor, s string) {
	if glyphs != nil {
		textPath(b, float64(x), float64(y), float64(size), fill, anchor, s)
		return
	}
	if anchor != "" {
		anchor = fmt.Sprintf(` text-anchor="%s"`, anchor)
	}
	fmt.Fprintf(b, `<text x="%d" y="%d" font-size="%d" fill="%s" font-family="monospace"%s>%s</text>`, x, y, size, fill, anchor, html.EscapeString(s))
}

func writeStarfield(b *strings.Builder) {
	stars := [][3]interface{}{
		{300, 14, colMuted}, {420, 10, colGold}, {530, 16, colMuted}, {640, 12, colMuted},
		{360, 292, colMuted}, {560, 296, colGold}, {620, 292, colMuted}, {680, 40, colMuted},
	}
	for _, s := range stars {
		rct(b, s[0].(int), s[1].(int), 2, 2, s[2].(string))
	}
}

// writeFrame draws the classic JRPG double window border with stepped pixel corners.
func writeFrame(b *strings.Builder, outer string) {
	x, y, w, h, c := 2, 2, cardW-4, cardH-4, 6
	rct(b, x+c, y, w-2*c, 3, outer)
	rct(b, x+c, y+h-3, w-2*c, 3, outer)
	rct(b, x, y+c, 3, h-2*c, outer)
	rct(b, x+w-3, y+c, 3, h-2*c, outer)
	for _, s := range [][4]int{{x, y, 1, 1}, {x + w, y, -1, 1}, {x, y + h, 1, -1}, {x + w, y + h, -1, -1}} {
		cx, cy, dx, dy := s[0], s[1], s[2], s[3]
		hx := cx
		if dx < 0 {
			hx = cx - 6
		}
		hy := cy + 3*dy
		if dy < 0 {
			hy = cy - 6
		}
		rct(b, hx, boolPick(dy > 0, cy+3, cy-6), 6, 3, outer)
		rct(b, boolPick(dx > 0, cx+3, cx-6), hy, 3, 6, outer)
	}
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="1"/>`, x+6, y+6, w-12, h-12, colFrameD)
}

func boolPick(cond bool, a, c int) int {
	if cond {
		return a
	}
	return c
}

func writeTorch(b *strings.Builder, x, y int, flicker bool, dur string) {
	rct(b, x, y+12, 4, 10, "#6b4a2b")
	rct(b, x-2, y+10, 8, 3, "#4a5568")
	anim := ""
	if flicker {
		anim = fmt.Sprintf(`<animate attributeName="opacity" values="1;0.5;1" dur="%s" repeatCount="indefinite"/>`, dur)
	}
	fmt.Fprintf(b, `<g>%s<rect x="%d" y="%d" width="6" height="6" fill="%s"/><rect x="%d" y="%d" width="4" height="6" fill="%s"/><rect x="%d" y="%d" width="2" height="5" fill="%s"/></g>`,
		anim, x-1, y+4, colFlameO, x, y, colFlameO, x+1, y+3, colFlameY)
}

// writeScene draws the dungeon diorama: brick wall, dithered floor, torches,
// and the class sprite.
func writeScene(b *strings.Builder, px, py, pw, ph int, class, borderColor string, flicker bool) {
	rct(b, px, py, pw, ph, colWall)
	row := 0
	for yy := py + 14; yy < py+ph-40; yy += 16 {
		rct(b, px, yy, pw, 1, colWallLt)
		off := 20
		if row%2 == 1 {
			off = 40
		}
		for xx := px + off; xx < px+pw-4; xx += 44 {
			rct(b, xx, yy-14, 1, 14, colWallLt)
		}
		row++
	}
	fy := py + ph - 36
	rct(b, px, fy, pw, 36, colFloorB)
	for r := 0; r < 6; r++ {
		for c := 0; c <= pw/6; c++ {
			if (r+c)%2 == 0 {
				rct(b, px+c*6, fy+r*6, 6, 6, colFloorA)
			}
		}
	}
	rct(b, px, fy, pw, 2, colShadow)
	writeTorch(b, px+22, py+26, flicker, "1.2s")
	writeTorch(b, px+pw-28, py+26, flicker, "1.6s")
	for _, g := range [][2]int{{px + 18, py + 22}, {px + 30, py + 30}, {px + pw - 32, py + 22}, {px + pw - 20, py + 30}} {
		rct(b, g[0], g[1], 2, 2, colFlameY)
	}
	sx := px + (pw-96)/2
	rct(b, sx+18, fy+8, 60, 6, colShadow)
	if sprite, ok := spriteB64[class]; ok {
		fmt.Fprintf(b, `<image x="%d" y="%d" width="96" height="96" href="data:image/png;base64,%s" image-rendering="pixelated" style="image-rendering:pixelated"/>`,
			sx, fy-82, sprite)
	}
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="1"/>`, px, py, pw, ph, borderColor)
}

func writeRibbon(b *strings.Builder, x, y, w, h int, accentDark, label string) {
	fmt.Fprintf(b, `<polygon points="%d,%d %d,%d %d,%d %d,%d %d,%d" fill="%s"/>`,
		x-10, y+3, x, y+3, x, y+h-3, x-10, y+h-3, x-4, y+h/2, colGoldD)
	fmt.Fprintf(b, `<polygon points="%d,%d %d,%d %d,%d %d,%d %d,%d" fill="%s"/>`,
		x+w+10, y+3, x+w, y+3, x+w, y+h-3, x+w+10, y+h-3, x+w+4, y+h/2, colGoldD)
	rct(b, x, y, w, h, accentDark)
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="1"/>`, x+1, y+1, w-2, h-2, colGold)
	size := 11
	if len(label) > 22 {
		size = 6
	} else if len(label) > 14 {
		size = 8
	}
	txt(b, x+w/2, y+h/2+4, size, colText, "middle", label)
}

func writeDitherBar(b *strings.Builder, x, y, w, h, pct int, color string) {
	rct(b, x, y, w, h, colBarBG)
	fill := w * pct / 100
	if fill > 0 {
		rct(b, x, y, fill, h, color)
		rct(b, x, y, fill, 2, "#ffffff22")
		if fill < w-4 {
			for i := 0; i < h/2; i++ {
				if i%2 == 0 {
					rct(b, x+fill, y+i*2, 2, 2, color)
				} else {
					rct(b, x+fill+2, y+i*2, 2, 2, color)
				}
			}
		}
	}
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#000000" stroke-opacity="0.45" stroke-width="1"/>`, x, y, w, h)
}

func writeGearSlot(b *strings.Builder, x, y, size int, item *equipment.Item) {
	rct(b, x, y, size, size, "#141729")
	fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="1"/>`, x, y, size, size, colFrameD)
	for _, c := range [][2]int{{x, y}, {x + size - 3, y}, {x, y + size - 3}, {x + size - 3, y + size - 3}} {
		rct(b, c[0], c[1], 3, 3, colGoldD)
	}
	if item == nil {
		return
	}
	if icon, ok := gearIcons[item.Slug]; ok {
		fmt.Fprintf(b, `<image x="%d" y="%d" width="%d" height="%d" href="%s" image-rendering="pixelated" style="image-rendering:pixelated"/>`,
			x+4, y+4, size-8, size-8, icon)
	}
}

// accentDark returns a dark shade of the accent for panel fills, derived by
// halving each RGB channel twice.
func accentDark(accent string) string {
	if len(accent) != 7 {
		return colBarBG
	}
	var r, g, bl int
	if _, err := fmt.Sscanf(accent[1:], "%02x%02x%02x", &r, &g, &bl); err != nil {
		return colBarBG
	}
	return fmt.Sprintf("#%02x%02x%02x", r/4+8, g/4+8, bl/4+16)
}

func render(login string, char *stats.Character, loadout equipment.Loadout) string {
	accent := ClassColor(char.Class)
	stage := cardStage(char.Level)
	tierColor := tierColors[stage]
	legendary := stage >= 3

	xpPct := 0
	if char.XPForLevel > 0 {
		xpPct = char.XPIntoLevel * 100 / char.XPForLevel
		if xpPct > 100 {
			xpPct = 100
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg width="840" height="370" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, cardW, cardH)
	rct(&b, 0, 0, cardW, cardH, colBG)
	writeStarfield(&b)
	frameColor := colBone
	if legendary {
		frameColor = colGold
	}
	writeFrame(&b, frameColor)

	// left: dungeon scene + nameplate ribbon
	px, py, pw, ph := 22, 22, 224, 216
	sceneBorder := colFrameD
	if stage >= 1 {
		sceneBorder = tierColor
	}
	writeScene(&b, px, py, pw, ph, char.Class, sceneBorder, legendary)
	writeRibbon(&b, px+22, py+ph+14, pw-44, 26, accentDark(accent), login)
	txt(&b, 28, cardH-16, 6, colFrameD, "", "GITHUB-RPG")

	// right: header
	rx := 270
	txt(&b, rx, 50, 20, colText, "", fmt.Sprintf("LV %d", char.Level))
	fmt.Fprintf(&b, `<rect x="%d" y="28" width="116" height="26" fill="none" stroke="%s" stroke-width="1"/>`, rx+124, accent)
	txt(&b, rx+182, 47, 9, accent, "middle", strings.ToUpper(char.Class))
	if stage >= 1 {
		txt(&b, cardW-28, 47, 10, tierColor, "end", "* "+tierLabels[stage])
	}
	txt(&b, rx, 74, 9, colMuted, "", strings.ToUpper(char.Title))
	rct(&b, rx, 86, cardW-28-rx, 2, colDivider)

	// right: stats
	y0, rowH := 96, 22
	barX := rx + 160
	barW := cardW - 28 - barX
	rows := []struct {
		name string
		val  int
	}{
		{"POWER", char.Strength}, {"INTELLECT", char.Intelligence}, {"WISDOM", char.Wisdom},
		{"AGILITY", char.Dexterity}, {"INFLUENCE", char.Charisma},
	}
	for i, row := range rows {
		y := y0 + i*rowH
		txt(&b, rx, y+11, 10, colMuted, "", row.name)
		if row.val >= 100 {
			txt(&b, barX-10, y+11, 10, colGold, "end", "MAX")
			writeDitherBar(&b, barX, y, barW, 13, 100, colGold)
		} else {
			txt(&b, barX-10, y+11, 11, colText, "end", fmt.Sprintf("%d", row.val))
			writeDitherBar(&b, barX, y, barW, 13, row.val, accent)
		}
	}

	// right: exp
	ey := y0 + 5*rowH + 6
	txt(&b, rx, ey+11, 10, colGold, "", "EXP")
	writeDitherBar(&b, rx+52, ey, cardW-28-(rx+52)-128, 13, xpPct, colGold)
	txt(&b, cardW-28, ey+11, 10, colText, "end", fmt.Sprintf("%d/%d", char.XPIntoLevel, char.XPForLevel))
	txt(&b, rx, ey+28, 8, colMuted, "", fmt.Sprintf("NEXT LV %d IN %d XP", char.Level+1, char.XPForLevel-char.XPIntoLevel))

	// right: gear slots
	gy := ey + 36
	txt(&b, rx, gy+18, 10, colMuted, "", "GEAR")
	gx, size, pitch := rx+22, 28, 126
	for i, item := range []*equipment.Item{loadout.Weapon, loadout.Shield, loadout.Accessory} {
		cx := gx + i*pitch + (pitch-size)/2
		writeGearSlot(&b, cx, gy, size, item)
		label, color := "EMPTY", colFrameD
		if item != nil {
			label, color = strings.ToUpper(item.Name), colMuted
		}
		txt(&b, gx+i*pitch+pitch/2, gy+size+13, 7, color, "middle", label)
	}

	b.WriteString("</svg>")
	return b.String()
}

// Card renders the RPG character card as an SVG string.
func Card(login string, char *stats.Character, loadout equipment.Loadout) (string, error) {
	return render(login, char, loadout), nil
}

var demoChars = map[string]*stats.Character{
	"Guardian":   {Class: "Guardian", Level: 15, Title: "The Architect", TotalXP: 14200, XPIntoLevel: 200, XPForLevel: 3400, Strength: 85, Intelligence: 45, Wisdom: 60, Dexterity: 50, Charisma: 55},
	"Berserker":  {Class: "Berserker", Level: 18, Title: "The Unstoppable", TotalXP: 22800, XPIntoLevel: 800, XPForLevel: 4200, Strength: 95, Intelligence: 30, Wisdom: 25, Dexterity: 75, Charisma: 40},
	"Paladin":    {Class: "Paladin", Level: 20, Title: "The Honorbound", TotalXP: 31000, XPIntoLevel: 1000, XPForLevel: 5000, Strength: 75, Intelligence: 55, Wisdom: 80, Dexterity: 45, Charisma: 85},
	"Rogue":      {Class: "Rogue", Level: 12, Title: "The Shadow", TotalXP: 9600, XPIntoLevel: 600, XPForLevel: 2800, Strength: 55, Intelligence: 70, Wisdom: 40, Dexterity: 95, Charisma: 65},
	"Sage":       {Class: "Sage", Level: 25, Title: "The Omniscient", TotalXP: 52000, XPIntoLevel: 2000, XPForLevel: 7500, Strength: 35, Intelligence: 95, Wisdom: 85, Dexterity: 45, Charisma: 60},
	"Knight":     {Class: "Knight", Level: 22, Title: "The Valiant", TotalXP: 38500, XPIntoLevel: 500, XPForLevel: 5800, Strength: 80, Intelligence: 50, Wisdom: 65, Dexterity: 60, Charisma: 70},
	"Battlemage": {Class: "Battlemage", Level: 17, Title: "The Spellblade", TotalXP: 19000, XPIntoLevel: 1000, XPForLevel: 3900, Strength: 65, Intelligence: 85, Wisdom: 55, Dexterity: 50, Charisma: 45},
	"Warlord":    {Class: "Warlord", Level: 52, Title: "The Conqueror", TotalXP: 121000, XPIntoLevel: 3000, XPForLevel: 9000, Strength: 90, Intelligence: 40, Wisdom: 50, Dexterity: 70, Charisma: 75},
	"Wanderer":   {Class: "Wanderer", Level: 14, Title: "The Pathfinder", TotalXP: 12000, XPIntoLevel: 400, XPForLevel: 3200, Strength: 60, Intelligence: 55, Wisdom: 70, Dexterity: 85, Charisma: 50},
}

// demoWeapons is each demo class's weapon slot, matching its language item.
var demoWeapons = map[string]string{
	"Guardian":   "go-compiler",
	"Berserker":  "rust-axe",
	"Paladin":    "typescript-lens",
	"Rogue":      "javascript-dagger",
	"Sage":       "python-flask",
	"Knight":     "csharp-bulwark",
	"Battlemage": "java-hammer",
	"Warlord":    "cpp-gauntlet",
	"Wanderer":   "go-compiler",
}

func demoLoadout(class string) equipment.Loadout {
	pick := func(slug string) *equipment.Item {
		if it, ok := equipment.Lookup(slug); ok {
			return &it
		}
		return nil
	}
	weapon := demoWeapons[class]
	shield := "csharp-bulwark"
	if weapon == shield {
		shield = "go-compiler"
	}
	accessory := "python-flask"
	if weapon == accessory {
		accessory = "typescript-lens"
	}
	return equipment.Loadout{Weapon: pick(weapon), Shield: pick(shield), Accessory: pick(accessory)}
}

// Demo renders a sample card for a class. A level > 0 overrides the sample
// character's level, so the evolution tiers can be previewed.
func Demo(class string, level int) (string, error) {
	char, ok := demoChars[class]
	if !ok {
		class = "Wanderer"
		char = demoChars[class]
	}
	if level > 0 {
		c := *char
		c.Level = level
		char = &c
	}
	return render(class, char, demoLoadout(class)), nil
}
