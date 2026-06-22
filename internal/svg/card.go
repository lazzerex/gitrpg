package svg

import (
	"bytes"
	"text/template"

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
	XPPct       int
	XPArcDash   int
	AccentColor string
	STR         int
	INT         int
	WIS         int
	DEX         int
	CHA         int
	STRBar      int
	INTBar      int
	WISBar      int
	DEXBar      int
	CHABar      int
	STRBar100   int
	INTBar100   int
	WISBar100   int
	DEXBar100   int
	CHABar100   int
	STRChartH   int
	INTChartH   int
	WISChartH   int
	DEXChartH   int
	CHAChartH   int
	STRChartY   int
	INTChartY   int
	WISChartY   int
	DEXChartY   int
	CHAChartY   int
}

func statBar(val, max int) int {
	if val <= 0 {
		return 0
	}
	b := val * max / 100
	if b > max {
		return max
	}
	return b
}

// classicSVG: 495×195. Stat numbers + mini bars.
const classicSVG = `<svg width="495" height="195" viewBox="0 0 495 195" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="495" height="195" fill="#0a001a"/>
  <rect x="0" y="0" width="495" height="195" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="489" height="189" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <text x="16" y="30" font-family="Courier New,Courier,monospace" font-size="16" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="16" y="36" width="52" height="16" fill="{{.AccentColor}}"/>
  <text x="42" y="48" font-family="Courier New,Courier,monospace" font-size="8" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="74" y="48" font-family="Courier New,Courier,monospace" font-size="9" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="16" y="64" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" font-style="italic">{{.Title}}</text>
  <line x1="10" y1="71" x2="485" y2="71" stroke="#3A1A7A" stroke-width="1"/>
  <text x="16" y="81" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">XP {{.TotalXP}}</text>
  <text x="479" y="81" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="end">{{.XPInto}} / {{.XPFor}} &#8594; LV{{.NextLevel}}</text>
  <rect x="15" y="85" width="465" height="7" fill="#180830"/>
  <rect x="15" y="85" width="{{.XPBarWidth}}" height="7" fill="{{.AccentColor}}"/>
  <rect x="15" y="85" width="465" height="7" fill="none" stroke="#3A1A7A" stroke-width="1"/>
  <line x1="10" y1="100" x2="485" y2="100" stroke="#3A1A7A" stroke-width="1"/>
  <rect x="108" y="100" width="1" height="56" fill="#3A1A7A"/>
  <rect x="201" y="100" width="1" height="56" fill="#3A1A7A"/>
  <rect x="294" y="100" width="1" height="56" fill="#3A1A7A"/>
  <rect x="387" y="100" width="1" height="56" fill="#3A1A7A"/>
  <text x="54"  y="114" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">POWER</text>
  <text x="154" y="114" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">INTELLECT</text>
  <text x="247" y="114" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">WISDOM</text>
  <text x="340" y="114" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">AGILITY</text>
  <text x="441" y="114" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">INFLUENCE</text>
  <text x="54"  y="138" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.STR}}</text>
  <text x="154" y="138" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.INT}}</text>
  <text x="247" y="138" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.WIS}}</text>
  <text x="340" y="138" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.DEX}}</text>
  <text x="441" y="138" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.CHA}}</text>
  <rect x="14"  y="145" width="80" height="4" fill="#180830"/>
  <rect x="114" y="145" width="80" height="4" fill="#180830"/>
  <rect x="207" y="145" width="80" height="4" fill="#180830"/>
  <rect x="300" y="145" width="80" height="4" fill="#180830"/>
  <rect x="401" y="145" width="80" height="4" fill="#180830"/>
  <rect x="14"  y="145" width="{{.STRBar}}" height="4" fill="{{.AccentColor}}"/>
  <rect x="114" y="145" width="{{.INTBar}}" height="4" fill="{{.AccentColor}}"/>
  <rect x="207" y="145" width="{{.WISBar}}" height="4" fill="{{.AccentColor}}"/>
  <rect x="300" y="145" width="{{.DEXBar}}" height="4" fill="{{.AccentColor}}"/>
  <rect x="401" y="145" width="{{.CHABar}}" height="4" fill="{{.AccentColor}}"/>
  <line x1="10" y1="157" x2="485" y2="157" stroke="#3A1A7A" stroke-width="1"/>
  <path transform="translate(10,165) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

// chartSVG: 495×210. Vertical bar chart stats section.
const chartSVG = `<svg width="495" height="210" viewBox="0 0 495 210" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="495" height="210" fill="#0a001a"/>
  <rect x="0" y="0" width="495" height="210" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="489" height="204" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <text x="16" y="30" font-family="Courier New,Courier,monospace" font-size="16" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="16" y="36" width="52" height="16" fill="{{.AccentColor}}"/>
  <text x="42" y="48" font-family="Courier New,Courier,monospace" font-size="8" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="74" y="48" font-family="Courier New,Courier,monospace" font-size="9" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="16" y="64" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" font-style="italic">{{.Title}}</text>
  <line x1="10" y1="71" x2="485" y2="71" stroke="#3A1A7A" stroke-width="1"/>
  <text x="16" y="81" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">XP {{.TotalXP}}</text>
  <text x="479" y="81" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="end">{{.XPInto}} / {{.XPFor}} &#8594; LV{{.NextLevel}}</text>
  <rect x="15" y="85" width="465" height="7" fill="#180830"/>
  <rect x="15" y="85" width="{{.XPBarWidth}}" height="7" fill="{{.AccentColor}}"/>
  <rect x="15" y="85" width="465" height="7" fill="none" stroke="#3A1A7A" stroke-width="1"/>
  <line x1="10" y1="100" x2="485" y2="100" stroke="#3A1A7A" stroke-width="1"/>
  <rect x="108" y="100" width="1" height="72" fill="#3A1A7A"/>
  <rect x="201" y="100" width="1" height="72" fill="#3A1A7A"/>
  <rect x="294" y="100" width="1" height="72" fill="#3A1A7A"/>
  <rect x="387" y="100" width="1" height="72" fill="#3A1A7A"/>
  <text x="54"  y="113" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">POWER</text>
  <text x="154" y="113" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">INTELLECT</text>
  <text x="247" y="113" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">WISDOM</text>
  <text x="340" y="113" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">AGILITY</text>
  <text x="441" y="113" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">INFLUENCE</text>
  <rect x="32"  y="117" width="44" height="55" fill="#180830"/>
  <rect x="124" y="117" width="44" height="55" fill="#180830"/>
  <rect x="217" y="117" width="44" height="55" fill="#180830"/>
  <rect x="310" y="117" width="44" height="55" fill="#180830"/>
  <rect x="402" y="117" width="44" height="55" fill="#180830"/>
  <rect x="32"  y="{{.STRChartY}}" width="44" height="{{.STRChartH}}" fill="{{.AccentColor}}"/>
  <rect x="124" y="{{.INTChartY}}" width="44" height="{{.INTChartH}}" fill="{{.AccentColor}}"/>
  <rect x="217" y="{{.WISChartY}}" width="44" height="{{.WISChartH}}" fill="{{.AccentColor}}"/>
  <rect x="310" y="{{.DEXChartY}}" width="44" height="{{.DEXChartH}}" fill="{{.AccentColor}}"/>
  <rect x="402" y="{{.CHAChartY}}" width="44" height="{{.CHAChartH}}" fill="{{.AccentColor}}"/>
  <line x1="10" y1="173" x2="485" y2="173" stroke="#3A1A7A" stroke-width="1"/>
  <text x="54"  y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.STR}}</text>
  <text x="154" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.INT}}</text>
  <text x="247" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.WIS}}</text>
  <text x="340" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.DEX}}</text>
  <text x="441" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.CHA}}</text>
  <path transform="translate(10,190) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

// statsSVG: 495×220. Left panel: stat list with Lucide-style fill icons + bars.
// Right panel: XP donut ring with level number. Vertical panel divider at x=270.
// Icon paths are 24×24 viewbox fill shapes scaled to 12px via scale(0.5).
//   POWER     — lightning bolt polygon
//   INTELLECT — hexagon (structured knowledge)
//   WISDOM    — five-pointed star
//   AGILITY   — right-pointing triangle
//   INFLUENCE — person silhouette (head + body)
const statsSVG = `<svg width="495" height="220" viewBox="0 0 495 220" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="495" height="220" fill="#0a001a"/>
  <rect x="0" y="0" width="495" height="220" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="489" height="214" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <text x="16" y="28" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="16" y="33" width="52" height="15" fill="{{.AccentColor}}"/>
  <text x="42" y="44" font-family="Courier New,Courier,monospace" font-size="8" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="74" y="44" font-family="Courier New,Courier,monospace" font-size="8" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="16" y="58" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" font-style="italic">{{.Title}}</text>
  <line x1="10" y1="66" x2="262" y2="66" stroke="#3A1A7A" stroke-width="1"/>
  <path transform="translate(8,72) scale(0.5)" fill="{{.AccentColor}}" d="M13 2L3 14H12L11 22L21 10H12L13 2Z"/>
  <text x="26" y="82" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">POWER</text>
  <text x="148" y="82" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.STR}}</text>
  <rect x="153" y="76" width="100" height="8" fill="#180830"/>
  <rect x="153" y="76" width="{{.STRBar100}}" height="8" fill="{{.AccentColor}}"/>
  <path transform="translate(8,92) scale(0.5)" fill="{{.AccentColor}}" d="M12 2L22 7V17L12 22L2 17V7L12 2Z"/>
  <text x="26" y="102" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">INTELLECT</text>
  <text x="148" y="102" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.INT}}</text>
  <rect x="153" y="96" width="100" height="8" fill="#180830"/>
  <rect x="153" y="96" width="{{.INTBar100}}" height="8" fill="{{.AccentColor}}"/>
  <path transform="translate(8,112) scale(0.5)" fill="{{.AccentColor}}" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
  <text x="26" y="122" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">WISDOM</text>
  <text x="148" y="122" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.WIS}}</text>
  <rect x="153" y="116" width="100" height="8" fill="#180830"/>
  <rect x="153" y="116" width="{{.WISBar100}}" height="8" fill="{{.AccentColor}}"/>
  <path transform="translate(8,132) scale(0.5)" fill="{{.AccentColor}}" d="M8 5L19 12L8 19Z"/>
  <text x="26" y="142" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">AGILITY</text>
  <text x="148" y="142" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.DEX}}</text>
  <rect x="153" y="136" width="100" height="8" fill="#180830"/>
  <rect x="153" y="136" width="{{.DEXBar100}}" height="8" fill="{{.AccentColor}}"/>
  <path transform="translate(8,152) scale(0.5)" fill="{{.AccentColor}}" d="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm0 2c-4.42 0-8 1.79-8 4v2h16v-2c0-2.21-3.58-4-8-4z"/>
  <text x="26" y="162" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD">INFLUENCE</text>
  <text x="148" y="162" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.CHA}}</text>
  <rect x="153" y="156" width="100" height="8" fill="#180830"/>
  <rect x="153" y="156" width="{{.CHABar100}}" height="8" fill="{{.AccentColor}}"/>
  <line x1="270" y1="66" x2="270" y2="180" stroke="#3A1A7A" stroke-width="1"/>
  <circle cx="380" cy="120" r="50" fill="none" stroke="#180830" stroke-width="10"/>
  <circle cx="380" cy="120" r="50" fill="none" stroke="{{.AccentColor}}" stroke-width="10" stroke-dasharray="{{.XPArcDash}} 314" transform="rotate(-90 380 120)"/>
  <text x="380" y="108" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">LVL</text>
  <text x="380" y="127" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.Level}}</text>
  <text x="380" y="143" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" text-anchor="middle">{{.XPPct}}%</text>
  <text x="380" y="185" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">{{.XPInto}} / {{.XPFor}} XP</text>
  <text x="380" y="197" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">&#8594; LV{{.NextLevel}}</text>
  <path transform="translate(10,205) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

// compactSVG: 300×130. Clean rect border.
const compactSVG = `<svg width="300" height="130" viewBox="0 0 300 130" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="300" height="130" fill="#0a001a"/>
  <rect x="0" y="0" width="300" height="130" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="294" height="124" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <text x="14" y="22" font-family="Courier New,Courier,monospace" font-size="12" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="14" y="27" width="44" height="14" fill="{{.AccentColor}}"/>
  <text x="36" y="38" font-family="Courier New,Courier,monospace" font-size="7" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="64" y="38" font-family="Courier New,Courier,monospace" font-size="8" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="14" y="51" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" font-style="italic">{{.Title}}</text>
  <rect x="14" y="57" width="272" height="6" fill="#180830"/>
  <rect x="14" y="57" width="{{.XPBarWidth}}" height="6" fill="{{.AccentColor}}"/>
  <rect x="14" y="57" width="272" height="6" fill="none" stroke="#3A1A7A" stroke-width="1"/>
  <text x="14"  y="72" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD">{{.TotalXP}} XP</text>
  <text x="286" y="72" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="end">LV{{.NextLevel}} in {{.XPFor}} XP</text>
  <line x1="14" y1="77" x2="286" y2="77" stroke="#3A1A7A" stroke-width="1"/>
  <rect x="63"  y="77" width="1" height="38" fill="#3A1A7A"/>
  <rect x="120" y="77" width="1" height="38" fill="#3A1A7A"/>
  <rect x="179" y="77" width="1" height="38" fill="#3A1A7A"/>
  <rect x="236" y="77" width="1" height="38" fill="#3A1A7A"/>
  <text x="31"  y="89" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="middle">POWER</text>
  <text x="91"  y="89" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="middle">INTELLECT</text>
  <text x="149" y="89" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="middle">WISDOM</text>
  <text x="207" y="89" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="middle">AGILITY</text>
  <text x="268" y="89" font-family="Courier New,Courier,monospace" font-size="6" fill="#AA88DD" text-anchor="middle">INFLUENCE</text>
  <text x="31"  y="110" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.STR}}</text>
  <text x="91"  y="110" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.INT}}</text>
  <text x="149" y="110" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.WIS}}</text>
  <text x="207" y="110" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.DEX}}</text>
  <text x="268" y="110" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.CHA}}</text>
  <path transform="translate(9,114) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

var (
	classicTmpl = template.Must(template.New("classic").Parse(classicSVG))
	chartTmpl   = template.Must(template.New("chart").Parse(chartSVG))
	statsTmpl   = template.Must(template.New("stats").Parse(statsSVG))
	compactTmpl = template.Must(template.New("compact").Parse(compactSVG))
)

const (
	chartMaxH     = 55
	chartBaseline = 172
	donutR        = 50
	donutCirc     = 314 // ≈ 2π*50
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
		XPPct:       xpPercent,
		XPArcDash:   xpPercent * donutCirc / 100,
		AccentColor: accent,
		STR:         char.Strength,
		INT:         char.Intelligence,
		WIS:         char.Wisdom,
		DEX:         char.Dexterity,
		CHA:         char.Charisma,
		STRBar:      statBar(char.Strength, 80),
		INTBar:      statBar(char.Intelligence, 80),
		WISBar:      statBar(char.Wisdom, 80),
		DEXBar:      statBar(char.Dexterity, 80),
		CHABar:      statBar(char.Charisma, 80),
		STRBar100:   statBar(char.Strength, 100),
		INTBar100:   statBar(char.Intelligence, 100),
		WISBar100:   statBar(char.Wisdom, 100),
		DEXBar100:   statBar(char.Dexterity, 100),
		CHABar100:   statBar(char.Charisma, 100),
		STRChartH:   statBar(char.Strength, chartMaxH),
		INTChartH:   statBar(char.Intelligence, chartMaxH),
		WISChartH:   statBar(char.Wisdom, chartMaxH),
		DEXChartH:   statBar(char.Dexterity, chartMaxH),
		CHAChartH:   statBar(char.Charisma, chartMaxH),
		STRChartY:   chartBaseline - statBar(char.Strength, chartMaxH),
		INTChartY:   chartBaseline - statBar(char.Intelligence, chartMaxH),
		WISChartY:   chartBaseline - statBar(char.Wisdom, chartMaxH),
		DEXChartY:   chartBaseline - statBar(char.Dexterity, chartMaxH),
		CHAChartY:   chartBaseline - statBar(char.Charisma, chartMaxH),
	}
}

func Card(login string, char *stats.Character, style string) (string, error) {
	data := buildData(login, char, 465)
	var buf bytes.Buffer
	tmpl := classicTmpl
	switch style {
	case "chart":
		tmpl = chartTmpl
	case "stats":
		tmpl = statsTmpl
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Compact(login string, char *stats.Character) (string, error) {
	var buf bytes.Buffer
	if err := compactTmpl.Execute(&buf, buildData(login, char, 272)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
