package svg

import (
	"bytes"
	"text/template"

	"github.com/lazzerex/gitrpg/internal/stats"
)

// ClassColor returns the hex accent color for the given class name.
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

type cardData struct {
	Login       string
	Level       int
	NextLevel   int
	Class       string
	Title       string
	TotalXP     int
	XPInto      int
	XPFor       int
	XPBarWidth  int
	AccentColor string
	STR         int
	INT         int
	WIS         int
	DEX         int
	CHA         int
}

// card dimensions: 495 × 195 (standard github stats card size)
const cardSVG = `<svg width="495" height="195" viewBox="0 0 495 195" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="0.5" y="0.5" rx="4.5" height="194" stroke="#30363d" width="494" fill="#0d1117"/>
  <text x="25" y="35" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="18" font-weight="bold" fill="#e6edf3">{{.Login}}</text>
  <text x="25" y="55" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="13" fill="{{.AccentColor}}">Level {{.Level}} {{.Class}}</text>
  <text x="25" y="71" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="11" fill="#6e7681" font-style="italic">{{.Title}}</text>
  <rect x="25" y="82" width="445" height="7" rx="3.5" fill="#21262d"/>
  <rect x="25" y="82" width="{{.XPBarWidth}}" height="7" rx="3.5" fill="{{.AccentColor}}"/>
  <text x="25" y="101" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681">{{.TotalXP}} XP</text>
  <text x="470" y="101" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="end">{{.XPInto}} / {{.XPFor}} to lv {{.NextLevel}}</text>
  <line x1="25" y1="112" x2="470" y2="112" stroke="#21262d" stroke-width="1"/>
  <text x="69" y="132" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="middle">STR</text>
  <text x="69" y="155" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="20" font-weight="bold" fill="#e6edf3" text-anchor="middle">{{.STR}}</text>
  <text x="158" y="132" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="middle">INT</text>
  <text x="158" y="155" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="20" font-weight="bold" fill="#e6edf3" text-anchor="middle">{{.INT}}</text>
  <text x="247" y="132" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="middle">WIS</text>
  <text x="247" y="155" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="20" font-weight="bold" fill="#e6edf3" text-anchor="middle">{{.WIS}}</text>
  <text x="336" y="132" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="middle">DEX</text>
  <text x="336" y="155" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="20" font-weight="bold" fill="#e6edf3" text-anchor="middle">{{.DEX}}</text>
  <text x="425" y="132" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" text-anchor="middle">CHA</text>
  <text x="425" y="155" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="20" font-weight="bold" fill="#e6edf3" text-anchor="middle">{{.CHA}}</text>
  <text x="25" y="183" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#484f58">gitrpg</text>
</svg>`

// compact dimensions: 300 × 120
const compactSVG = `<svg width="300" height="120" viewBox="0 0 300 120" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect x="0.5" y="0.5" rx="4.5" height="119" stroke="#30363d" width="299" fill="#0d1117"/>
  <text x="20" y="28" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="15" font-weight="bold" fill="#e6edf3">{{.Login}}</text>
  <text x="20" y="46" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="11" fill="{{.AccentColor}}">Level {{.Level}} {{.Class}}</text>
  <text x="20" y="60" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="10" fill="#6e7681" font-style="italic">{{.Title}}</text>
  <rect x="20" y="70" width="260" height="5" rx="2.5" fill="#21262d"/>
  <rect x="20" y="70" width="{{.XPBarWidth}}" height="5" rx="2.5" fill="{{.AccentColor}}"/>
  <text x="20" y="84" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681">{{.TotalXP}} XP · lv {{.NextLevel}} in {{.XPFor}} XP</text>
  <line x1="20" y1="92" x2="280" y2="92" stroke="#21262d" stroke-width="1"/>
  <text x="40" y="107" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681" text-anchor="middle">STR {{.STR}}</text>
  <text x="92" y="107" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681" text-anchor="middle">INT {{.INT}}</text>
  <text x="150" y="107" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681" text-anchor="middle">WIS {{.WIS}}</text>
  <text x="208" y="107" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681" text-anchor="middle">DEX {{.DEX}}</text>
  <text x="260" y="107" font-family="'Segoe UI',Ubuntu,Sans-Serif" font-size="9" fill="#6e7681" text-anchor="middle">CHA {{.CHA}}</text>
</svg>`

var (
	cardTmpl    = template.Must(template.New("card").Parse(cardSVG))
	compactTmpl = template.Must(template.New("compact").Parse(compactSVG))
)

func buildData(login string, char *stats.Character, barMaxWidth int) cardData {
	xpPercent := 0
	if char.XPForLevel > 0 {
		xpPercent = char.XPIntoLevel * 100 / char.XPForLevel
		if xpPercent > 100 {
			xpPercent = 100
		}
	}
	accent, ok := classColors[char.Class]
	if !ok {
		accent = classColors["Wanderer"]
	}
	return cardData{
		Login:       login,
		Level:       char.Level,
		NextLevel:   char.Level + 1,
		Class:       char.Class,
		Title:       char.Title,
		TotalXP:     char.TotalXP,
		XPInto:      char.XPIntoLevel,
		XPFor:       char.XPForLevel,
		XPBarWidth:  xpPercent * barMaxWidth / 100,
		AccentColor: accent,
		STR:         char.Strength,
		INT:         char.Intelligence,
		WIS:         char.Wisdom,
		DEX:         char.Dexterity,
		CHA:         char.Charisma,
	}
}

func Card(login string, char *stats.Character) (string, error) {
	var buf bytes.Buffer
	if err := cardTmpl.Execute(&buf, buildData(login, char, 445)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Compact(login string, char *stats.Character) (string, error) {
	var buf bytes.Buffer
	if err := compactTmpl.Execute(&buf, buildData(login, char, 260)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
