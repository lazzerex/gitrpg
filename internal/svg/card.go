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

// 32×32 pixel art icons from Kyrise's RPG Icon Pack, embedded as base64 PNGs.
const (
	iconSword     = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABLklEQVRYR9WWMQ6CMBSGYdRbOABx4gxM3oTdWTa5gd6EhOhE4uLMiCwOnMJF428e4ZFSQCmtLH/6StLv/f1bsC3Lelrix+6oT1p+L2IGQFEUrDPP82is1InaAe0ASZIwB1zXxVi1E7UDOgGoc4QxjmOMfd+HqnaiGTDtAFqcEB2xWZ2QnfFZQIwGmCUTQ65ZpVsxBOBbJ9ofOeFafwUgdSLPc8xHUQQNwxCaZRm0LEvh13WMA8YACEEWmx3qx9MNur5uWedTOmAMAANZ7T9/VI/LAbq8n6V7LwwEFUcqjpsOACxMaa+qqt0x60NZBrQDOI4j7JTqQRBgnhxK05Rt/zf3AAufTgAG0rCBmkJGCJCUskD6iwPGAPSdWuYEvTylA8YDdG0V6lNkoM8BKcAL+MPCG+UmDIIAAAAASUVORK5CYIIA"
	iconSword2    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABUklEQVRYR82WO26DQBCGQ5sz0LgAlIpbuOAguE5PZ05BR02PlN5SmjSULFKUgjqda1v5pVE8y7I2j/VAM+yygm++2Qfey/91ubn/u/W0tpPm7Ue2AdC2LTKNoogydmpiYGAzAJS+axOjBsQBlFJgCMPQ6ZwYNbAZACqFKxN3DYgBVFWFb8dxzHbAtU2MGpAAoEyxJdOG5BrEtM2KA6xlQj/cjGeL7aBZamIxwCQTTdNgfJZliGmastVTFMVkA5sBeAiEDJRlifFBECB2Xcei/qc15WfDOieeAcBM5HnOavy1O6DdfP8ivn2+I9Z1bay9tZO9ediACUkAZoJq7e2P6D+rE+Lrz4e19ksMiAOgBLTO+75nRaJZT51kaGwuTFkFLHNxAMqMqPTMkyTBI9/3EcmUbmK2AXEAw1LVk8FcsYBi/GwDkgB39qnBY+uxPMfAqgBX+oXuIQYvYaIAAAAASUVORK5CYIIA"
	iconSpellbook = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABN0lEQVRYR2NkIABYWFj+E1JDjPyfP38YsanDKoiscMAcALPYzc0N7B4ZGRliPIqh5smTJ2CxXbt2gWn0kMAZAgPmAJjFioqKYBc7OjqC6f3795MVAjD98+bNIy4EBr0DtBuOgX1ytcEKTBPiUz0ECFmILk91B6AnBHQL0eWHnwNgPob51FZNBMycsuMWiudZFvmg5CKq5YIBdwB6HP+J2wIWgvmY7mmA7g5At5AQn+q5gJCF6PJUdwCjWzM4mv/vqgXThPijDhh+IUBqo4DqIUB3B6iqqoLtpFubEL1FNGAOgLWGy8y/gENAQV2Z1NAHq3etPQKm79+/T1qbcNA4YFqcONjlghIqZIVA6bKHYH0ktwdgITBgDkDvF5DlfQYGBpJ7Rrj6BXRzAMwianVKYeaR3DumlwMAH1WIMKq8RIAAAAAASUVORK5CYIIA"
	iconStaff     = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAA9ElEQVRYR2NkYGD4z4AfMBKQp0gaZPjgcEDvaYg7puy4BaZ/HZoMpp/ungLzIU1CAh4CA+4AxeabKHFJ9xAYcAdIu+aghIDE9Xtg/rufj8D0/ddXaJIW4GlgIB0A8xlKdlQU1QGLuwZageknT56A6W3btlE1JJCz1oA7AL1EAztIVVUVLO7o6EiTkMBXuAy4A1DSBq1CgpjilaYhQYwDaBoSQ8oBNAkJUkJg0DiAqg4hJwQGjQOo4hBKQmDQOIAih1AjBAaNA8hyCDVDYNA4gCSH0CIEBo0D8Dpk1qxZYHlahsCgcQDWVjdVm9borVkcfKy9cADoJMwbtiyB5QAAAABJRU5ErkJgggAA"
	iconBow       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABOUlEQVRYR2NkwA3+o0kx4lFLthQ+QwfMAWCLZ6c7gX2VOnMfzHd0C4EBcwCKxRuPXwf7fMul5+jxS9WQQDZscDig0E0d7OP+XTdRfK4dygvmX139mappAiMEBqMDYI4ER5F3jww4BLaWPKFKSBATAoPGATAfg0Oi9LoOmN+teYWikCAlBAaNA6jqEHJCgLYOgNUBWEpCXCUgRWmClJKQvg64dv8peolIqA4gKyRwhsBAOAAlccGK5NsvPqHXikSFBLElJjbDwEE5kA6gakgQqkUJtgkpDQlKHEBsSBCqC9AbtyjtDEIJCqSYUJqguQNQQgJWUsKyKZZcgt6GBPNhUQmThLW4iAmBQeMAFIf46Eli9SkhQaRWNtjzpITAoHEAikMI+RiLPIqnyQmBQeMAMjyPqYWSEKCKAwB3aeIdjAqGWAAAAABJRU5ErkJgggAA"
	iconKey       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAA7klEQVRYR+2WwQ3DIAxFyTIwDWuwDQuwBtPAMq2E5IMtN8VAMVKTC5IFyfPnx/Zlxp/Xh6OX5JWizeTFagDtw6UUNlHnHMS7kuvaxGWuAYAyr7U2Lu99W3POiBPixpjbJCUKnAXA3HUDBCV+roAaAGQIa0oJeYDGl3tAEwBcju7aWtviMca2giLfMhcVC64OgBIaAEiJEMJQ5jMKnAUg/e9pA5FUQnp2qPA8ACsUQIOIhgfOAICBhM4FvRVwpg7czgU7ANie0Nv/V5hQDYAdw3f2gjMAqPvhPnZ4gHX/dgDq4r/yANuOaXB7IZoFeAN+/6whBbi1kAAAAABJRU5ErkJgggAA"
	iconGem       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABRklEQVRYR9WXMQrCUAxAFXXp4iAIXXX2AN5AwSN08wJOHsJDuLSewd3FzcXZVa1bh0qhiELsB/tNfxIHoy6B337z8pL+0npN+VdXzl/7W4B7hTlxQeINRWI1AEh87PWAwz8cIC59H+L0fDZi2IWxb3ytXAOgVHl3OASeRhRBPPX7ENfXq9gE14AaAJrYNNo2YNYlJigDagDOxJQBiYkqA2oAvMSj0bPIIIBwK+Jlu0UPSNdM2AZ+A+AehryXpGWA2tRcrd5OStSAOsCs0wHSQauFFjX2vNK66fE+z9H7N1kG67s05RnQBDCEMIxckEWSiCt/U2H9gzoAy4SZBduAq+e2Kta7oKoV3wBwmrABJJVTM2CbQmfimwCoiXm7DeuTOK58zqnTkZoBpwkNgJIJpDppQR9/Gal9F1AtFV8XKxNnIDaoAzwAtSfaIXzqsLAAAAAASUVORK5CYIIA"
	iconPotion    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABIElEQVRYR2NkIB78J14pWCUjMeqJUgQ1aHA4wEdPEuwef0tNFA9uPH4dzN9y6TlMnCjPEaUIOQQGwgEoQU+CA4gKCWJCYHA4YOVWiDu4/m0B0z4+PihpYMsWiPg3Joh4uDfcb3g9SXQIDBoHHD75CuxDW3MxlBBAF6dZCAyYA2DehUVFmBdqObdqG4SP5HPa5IKBcADMJ+BsQIIDiEngxFUYyCXhoHHAx6VXUBIBf7QOSfmfqASCmswYUKJgwB1wtfYwivu0m23pGwKjDhg0IXCTcTs47oOa2gYmDQyYA9bVVaHkArqHwKgDBl0IiHoXgNPElCZxkop5oqpMbLUhegjQzQE5dS/R6ihULs1DYCAdgNIywhMMpEQrSS0imjgAAPoJziGHZTDaAAAAAElFTkSuQmCC"
	iconShield    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABmElEQVRYR2NkGGDASKL9/4lUT7S5RCuEWkx3B6BY2Hsau/2Hb71BCZgN0aLoAYXTo4RCYMAcALYY5mOYD03uzwT77NGjR1iTgpycHFj8jGI6rhDB8DCuEBgwB4AtDlj6GuwDWzURMD0jSg1Mq6qqovhMRkYGzH/y5AmK+O3bt8F87YZjBEMCPQQGlwOuNlih+IDUEIBphoUEUu6AexxvCNDTAeCgV8teAXZ0nOQdMH3sGGocwuI2Pj4eJTfAUv/ChQuxphUrK0hILnquAqZvTY2ABQ4jLAQGlwMcfu/DmrphPoH5FBa3sKiChQx6yMFyywFWJ+JCYNA4ACULIHFe2beilBMzF6wE82EOx6WP5BAYcAegl3DoaYDRrRnsxv+7asE03dMAzRwA8xF6yUe3cmAgHABLbygFEswhMEly6wJYVCGXgPCiEC2lDy4HwIMFmspJDQGYz2HmkBwCA+EArGkBJgirJQm1CWG1Hj6f40oDg8YBKA6BcWDtBVxFNBYfE/IoA0n9goFwALpn6d41G3QOIBT9JMsDANCtYDCjzkYBAAAAAElFTkSuQmCC"
	iconScroll    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAB30lEQVRYR2NkwA3+45BixKOHZCl8hg2YA8AWP1iRiNU3ChHzYeJUCQlshgyYA1As/vrmA9in3CICYBrGh3lfO2c9VUICOQQGlwNg3sMVEjQPAXo6ABz0BxrtwHYqqCujxDl6GoDxe5adBaubsukSRWkBlAYGhwNg+f7BzbtgH4kKC4Pp12/fooQMzOfr7z1EKSeeXvlIVkjAQ2DAHfDxQCvYB+cP7sRaAq45CykXYD43SOAF8x+chPsczL+6+jNJIQEPgYF0AMzF4MSI7hBiff7hOhNZaQKjJBxIB2ANCa2cLrA4rjiH+TxQSR4lBEKMIXWIQ/0hvGkCZ20IC4mBcABKSHj3yGBN7eg+L4kyBqsTlFBByU2wkhWpHYESIgRbRAPuAGkdfrCLBTT/gWl0n8PiGuZTmPdgIQHjv39xBxKS0JIWljYIhsBAOABcHqAHPXo+XxqqD/YRrM6A1ZIwH8JCBFZ3wNIIMSEwOBygHQop62EAVsbDouRoTRCKPC6fwxThKhdwlgODxgFYajeUOgMlGBgYGCbPXAEWstEQBNOwtIAe9wwMDGDPEwyBwegAlJISxmmJ1EXxOXrI4KoTyAkBmjsA3QJi+4C4OrMoZT96yBAsCXGkE3RzQHyyHAAA1U1oLEglgK4AAAAASUVORK5CYIIA"
	iconHelmet    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAErSURBVFhH7ZTRDcIwDES7BAPwwx9TsAB7MAbfjIHELkzADAwB2NKTnGuiEKCkSJx0H7HP51ObdvjjBdwqnBy5pZEfR7IAHM4pgeoffBuJIfhGADdY7i/O1e7k3B6vTg1AHR1z+JhhK7oF8AEMIYawtY+vLahhHgGAGumrqNUBvraghm4BXMCl4qzEOF5AI/XcjBFdqI1QEiacPIAaLdabhNTj8mhc0mcCjpAVlgwnC8BnpMbKuPwZHb6hN4I3ZhMAYsAj1HOpzln92GMLFfMIwOXRQYy1rizp9FLaQkX3AMAF/FoxwJg6OtVrgIy+isSwWwB9FRhrH1LXANp/sAoXMtg9AI8QQ/olomPuJwMAH9Ag1E0g8PonFgMf7BkARJNnjFr1VbQaNuiH4Q5NHiqUMm6nXQAAAABJRU5ErkJggg=="
)

var classIconMap = map[string]string{
	"Guardian":   iconShield,
	"Knight":     iconSword,
	"Paladin":    iconScroll,
	"Berserker":  iconSword2,
	"Warlord":    iconHelmet,
	"Sage":       iconSpellbook,
	"Battlemage": iconStaff,
	"Rogue":      iconKey,
	"Wanderer":   iconBow,
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
	ClassIcon   string
	PotionIcon  string
	GemIcon     string
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
  <image x="455" y="4" width="36" height="36" href="{{.ClassIcon}}" image-rendering="pixelated"/>
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
  <image x="455" y="4" width="36" height="36" href="{{.ClassIcon}}" image-rendering="pixelated"/>
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
  <rect x="132" y="117" width="44" height="55" fill="#180830"/>
  <rect x="225" y="117" width="44" height="55" fill="#180830"/>
  <rect x="318" y="117" width="44" height="55" fill="#180830"/>
  <rect x="419" y="117" width="44" height="55" fill="#180830"/>
  <rect x="32"  y="{{.STRChartY}}" width="44" height="{{.STRChartH}}" fill="{{.AccentColor}}"/>
  <rect x="132" y="{{.INTChartY}}" width="44" height="{{.INTChartH}}" fill="{{.AccentColor}}"/>
  <rect x="225" y="{{.WISChartY}}" width="44" height="{{.WISChartH}}" fill="{{.AccentColor}}"/>
  <rect x="318" y="{{.DEXChartY}}" width="44" height="{{.DEXChartH}}" fill="{{.AccentColor}}"/>
  <rect x="419" y="{{.CHAChartY}}" width="44" height="{{.CHAChartH}}" fill="{{.AccentColor}}"/>
  <line x1="10" y1="173" x2="485" y2="173" stroke="#3A1A7A" stroke-width="1"/>
  <text x="54"  y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.STR}}</text>
  <text x="154" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.INT}}</text>
  <text x="247" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.WIS}}</text>
  <text x="340" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.DEX}}</text>
  <text x="441" y="186" font-family="Courier New,Courier,monospace" font-size="10" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.CHA}}</text>
  <path transform="translate(10,190) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

// statsSVG: 495×220. Left panel: stat list with icon paths + bars. Right: XP donut.
const statsSVG = `<svg width="495" height="220" viewBox="0 0 495 220" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="495" height="220" fill="#0a001a"/>
  <rect x="0" y="0" width="495" height="220" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="489" height="214" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <image x="455" y="4" width="36" height="36" href="{{.ClassIcon}}" image-rendering="pixelated"/>
  <text x="16" y="28" font-family="Courier New,Courier,monospace" font-size="14" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="16" y="33" width="52" height="15" fill="{{.AccentColor}}"/>
  <text x="42" y="44" font-family="Courier New,Courier,monospace" font-size="8" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="74" y="44" font-family="Courier New,Courier,monospace" font-size="8" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="16" y="58" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" font-style="italic">{{.Title}}</text>
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
  <circle cx="380" cy="120" r="50" fill="none" stroke="#180830" stroke-width="10"/>
  <circle cx="380" cy="120" r="50" fill="none" stroke="{{.AccentColor}}" stroke-width="10" stroke-dasharray="{{.XPArcDash}} 314" stroke-linecap="round" transform="rotate(-90 380 120)"/>
  <text x="380" y="103" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" text-anchor="middle">LVL</text>
  <text x="380" y="122" font-family="Courier New,Courier,monospace" font-size="22" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.Level}}</text>
  <text x="380" y="137" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">{{.TotalXP}} XP</text>
  <text x="380" y="150" font-family="Courier New,Courier,monospace" font-size="8" fill="{{.AccentColor}}" text-anchor="middle">{{.XPPct}}%</text>
  <text x="380" y="186" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">{{.XPInto}} / {{.XPFor}} XP</text>
  <text x="380" y="198" font-family="Courier New,Courier,monospace" font-size="7" fill="#AA88DD" text-anchor="middle">&#8594; LV{{.NextLevel}}</text>
  <path transform="translate(10,205) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
</svg>`

// compactSVG: 300×130.
const compactSVG = `<svg width="300" height="130" viewBox="0 0 300 130" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="300" height="130" fill="#0a001a"/>
  <rect x="0" y="0" width="300" height="130" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="3" y="3" width="294" height="124" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <image x="277" y="3" width="22" height="22" href="{{.ClassIcon}}" image-rendering="pixelated"/>
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
	arcDash := xpPercent * donutCirc / 100
	if xpPercent > 0 && arcDash < 20 {
		arcDash = 20
	}
	accent, ok := classColors[char.Class]
	if !ok {
		accent = classColors["Wanderer"]
	}
	icon, ok := classIconMap[char.Class]
	if !ok {
		icon = iconGem
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
		XPArcDash:   arcDash,
		AccentColor: accent,
		ClassIcon:   icon,
		PotionIcon:  iconPotion,
		GemIcon:     iconGem,
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

var demoChars = map[string]*stats.Character{
	"Guardian":   {Class: "Guardian", Level: 15, Title: "The Architect", TotalXP: 14200, XPIntoLevel: 200, XPForLevel: 3400, Strength: 85, Intelligence: 45, Wisdom: 60, Dexterity: 50, Charisma: 55},
	"Berserker":  {Class: "Berserker", Level: 18, Title: "The Unstoppable", TotalXP: 22800, XPIntoLevel: 800, XPForLevel: 4200, Strength: 95, Intelligence: 30, Wisdom: 25, Dexterity: 75, Charisma: 40},
	"Paladin":    {Class: "Paladin", Level: 20, Title: "The Honorbound", TotalXP: 31000, XPIntoLevel: 1000, XPForLevel: 5000, Strength: 75, Intelligence: 55, Wisdom: 80, Dexterity: 45, Charisma: 85},
	"Rogue":      {Class: "Rogue", Level: 12, Title: "The Shadow", TotalXP: 9600, XPIntoLevel: 600, XPForLevel: 2800, Strength: 55, Intelligence: 70, Wisdom: 40, Dexterity: 95, Charisma: 65},
	"Sage":       {Class: "Sage", Level: 25, Title: "The Omniscient", TotalXP: 52000, XPIntoLevel: 2000, XPForLevel: 7500, Strength: 35, Intelligence: 95, Wisdom: 85, Dexterity: 45, Charisma: 60},
	"Knight":     {Class: "Knight", Level: 22, Title: "The Valiant", TotalXP: 38500, XPIntoLevel: 500, XPForLevel: 5800, Strength: 80, Intelligence: 50, Wisdom: 65, Dexterity: 60, Charisma: 70},
	"Battlemage": {Class: "Battlemage", Level: 17, Title: "The Spellblade", TotalXP: 19000, XPIntoLevel: 1000, XPForLevel: 3900, Strength: 65, Intelligence: 85, Wisdom: 55, Dexterity: 50, Charisma: 45},
	"Warlord":    {Class: "Warlord", Level: 30, Title: "The Conqueror", TotalXP: 78000, XPIntoLevel: 3000, XPForLevel: 9000, Strength: 90, Intelligence: 40, Wisdom: 50, Dexterity: 70, Charisma: 75},
	"Wanderer":   {Class: "Wanderer", Level: 14, Title: "The Pathfinder", TotalXP: 12000, XPIntoLevel: 400, XPForLevel: 3200, Strength: 60, Intelligence: 55, Wisdom: 70, Dexterity: 85, Charisma: 50},
}

func Demo(class string) (string, error) {
	char, ok := demoChars[class]
	if !ok {
		char = demoChars["Wanderer"]
	}
	var buf bytes.Buffer
	if err := classicTmpl.Execute(&buf, buildData(class, char, 465)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
