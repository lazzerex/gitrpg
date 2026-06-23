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

// 32×32 pixel art class icons (base64 PNG).
const (
	iconSword     = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABLklEQVRYR9WWMQ6CMBSGYdRbOABx4gxM3oTdWTa5gd6EhOhE4uLMiCwOnMJF428e4ZFSQCmtLH/6StLv/f1bsC3Lelrix+6oT1p+L2IGQFEUrDPP82is1InaAe0ASZIwB1zXxVi1E7UDOgGoc4QxjmOMfd+HqnaiGTDtAFqcEB2xWZ2QnfFZQIwGmCUTQ65ZpVsxBOBbJ9ofOeFafwUgdSLPc8xHUQQNwxCaZRm0LEvh13WMA8YACEEWmx3qx9MNur5uWedTOmAMAANZ7T9/VI/LAbq8n6V7LwwEFUcqjpsOACxMaa+qqt0x60NZBrQDOI4j7JTqQRBgnhxK05Rt/zf3AAufTgAG0rCBmkJGCJCUskD6iwPGAPSdWuYEvTylA8YDdG0V6lNkoM8BKcAL+MPCG+UmDIIAAAAASUVORK5CYIIA"
	iconSword2    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABUklEQVRYR82WO26DQBCGQ5sz0LgAlIpbuOAguE5PZ05BR02PlN5SmjSULFKUgjqda1v5pVE8y7I2j/VAM+yygm++2Qfey/91ubn/u/W0tpPm7Ue2AdC2LTKNoogydmpiYGAzAJS+axOjBsQBlFJgCMPQ6ZwYNbAZACqFKxN3DYgBVFWFb8dxzHbAtU2MGpAAoEyxJdOG5BrEtM2KA6xlQj/cjGeL7aBZamIxwCQTTdNgfJZliGmastVTFMVkA5sBeAiEDJRlifFBECB2Xcei/qc15WfDOieeAcBM5HnOavy1O6DdfP8ivn2+I9Z1bay9tZO9ediACUkAZoJq7e2P6D+rE+Lrz4e19ksMiAOgBLTO+75nRaJZT51kaGwuTFkFLHNxAMqMqPTMkyTBI9/3EcmUbmK2AXEAw1LVk8FcsYBi/GwDkgB39qnBY+uxPMfAqgBX+oXuIQYvYaIAAAAASUVORK5CYIIA"
	iconSpellbook = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABN0lEQVRYR2NkIABYWFj+E1JDjPyfP38YsanDKoiscMAcALPYzc0N7B4ZGRliPIqh5smTJ2CxXbt2gWn0kMAZAgPmAJjFioqKYBc7OjqC6f3795MVAjD98+bNIy4EBr0DtBuOgX1ytcEKTBPiUz0ECFmILk91B6AnBHQL0eWHnwNgPob51FZNBMycsuMWiudZFvmg5CKq5YIBdwB6HP+J2wIWgvmY7mmA7g5At5AQn+q5gJCF6PJUdwCjWzM4mv/vqgXThPijDhh+IUBqo4DqIUB3B6iqqoLtpFubEL1FNGAOgLWGy8y/gENAQV2Z1NAHq3etPQKm79+/T1qbcNA4YFqcONjlghIqZIVA6bKHYH0ktwdgITBgDkDvF5DlfQYGBpJ7Rrj6BXRzAMwianVKYeaR3DumlwMAH1WIMKq8RIAAAAAASUVORK5CYIIA"
	iconStaff     = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAA9ElEQVRYR2NkYGD4z4AfMBKQp0gaZPjgcEDvaYg7puy4BaZ/HZoMpp/ungLzIU1CAh4CA+4AxeabKHFJ9xAYcAdIu+aghIDE9Xtg/rufj8D0/ddXaJIW4GlgIB0A8xlKdlQU1QGLuwZageknT56A6W3btlE1JJCz1oA7AL1EAztIVVUVLO7o6EiTkMBXuAy4A1DSBq1CgpjilaYhQYwDaBoSQ8oBNAkJUkJg0DiAqg4hJwQGjQOo4hBKQmDQOIAih1AjBAaNA8hyCDVDYNA4gCSH0CIEBo0D8Dpk1qxZYHlahsCgcQDWVjdVm9borVkcfKy9cADoJMwbtiyB5QAAAABJRU5ErkJgggAA"
	iconBow       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABOUlEQVRYR2NkwA3+o0kx4lFLthQ+QwfMAWCLZ6c7gX2VOnMfzHd0C4EBcwCKxRuPXwf7fMul5+jxS9WQQDZscDig0E0d7OP+XTdRfK4dygvmX139appAiMEBqMDYI4ER5F3jww4BLaWPKFKSBATAoPGATAfg0Oi9LoOmN+teYWikCAlBAaNA6jqEHJCgLYOgNUBWEpCXCUgRWmClJKQvg64dv8peolIqA4gKyRwhsBAOAAlccGK5NsvPqHXikSFBLElJjbDwEE5kA6gakgQqkUJtgkpDQlKHEBsSBCqC9AbtyjtDEIJCqSYUJqguQNQQgJWUsKyKZZcgt6GBPNhUQmThLW4iAmBQeMAFIf46Eli9SkhQaRWNtjzpITAoHEAikMI+RiLPIqnyQmBQeMAMjyPqYWSEKCKAwB3aeIdjAqGWAAAAABJRU5ErkJgggAA"
	iconKey       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAA7klEQVRYR+2WwQ3DIAxFyTIwDWuwDQuwBtPAMq2E5IMtN8VAMVKTC5IFyfPnx/Zlxp/Xh6OX5JWizeTFagDtw6UUNlHnHMS7kuvaxGWuAYAyr7U2Lu99W3POiBPixpjbJCUKnAXA3HUDBCV+roAaAGQIa0oJeYDGl3tAEwBcju7aWtviMca2giLfMhcVC64OgBIaAEiJEMJQ5jMKnAUg/e9pA5FUQnp2qPA8ACsUQIOIhgfOAICBhM4FvRVwpg7czgU7ANie0Nv/V5hQDYAdw3f2gjMAqPvhPnZ4gHX/dgDq4r/yANuOaXB7IZoFeAN+/6whBbi1kAAAAABJRU5ErkJgggAA"
	iconGem       = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABRklEQVRYR9WXMQrCUAxAFXXp4iAIXXX2AN5AwSN08wJOHsJDuLSewd3FzcXZVa1bh0qhiELsB/tNfxIHoy6B337z8pL+0npN+VdXzl/7W4B7hTlxQeINRWI1AEh87PWAwz8cIC59H+L0fDZi2IWxb3ytXAOgVHl3OASeRhRBPPX7ENfXq9gE14AaAJrYNNo2YNYlJigDagDOxJQBiYkqA2oAvMSj0bPIIIBwK+Jlu0UPSNdM2AZ+A+AehryXpGWA2tRcrd5OStSAOsCs0wHSQauFFjX2vNK66fE+z9H7N1kG67s05RnQBDCEMIxckEWSiCt/U2H9gzoAy4SZBduAq+e2Kta7oKoV3wBwmrABJJVTM2CbQmfimwCoiXm7DeuTOK58zqnTkZoBpwkNgJIJpDppQR9/Gal9F1AtFV8XKxNnIDaoAzwAtSfaIXzqsLAAAAAASUVORK5CYIIA"
	iconPotion    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABIElEQVRYR2NkIB78J14pWCUjMeqJUgQ1aHA4wEdPEuwef0tNFA9uPH4dzN9y6TlMnCjPEaUIOQQGwgEoQU+CA4gKCWJCYHA4YOVWiDu4/m0B0z4+PihpYMsWiPg3Joh4uDfcb3g9SXQIDBoHHD75CuxDW3MxlBBAF6dZCAyYA2DehUVFmBdqObdqG4SP5HPa5IKBcADMJ+BsQIIDiEngxFUYyCXhoHHAx6VXUBIBf7QOSfmfqASCmswYUKJgwB1wtfYwivu0m23pGwKjDhg0IXCTcTs47oOa2gYmDQyYA9bVVaHkArqHwKgDBl0IiHoXgNPElCZxkop5oqpMbLUhegjQzQE5dS/R6ihULs1DYCAdgNIywhMMpEQrSS0imjgAAPoJziGHZTDaAAAAAElFTkSuQmCC"
	iconShield    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAABmElEQVRYR2NkGGDASKL9/4lUT7S5RCuEWkx3B6BY2Hsau/2Hb71BCZgN0aLoAYXTo4RCYMAcALYY5mOYD03uzwT77NGjR1iTgpycHFj8jGI6rhDB8DCuEBgwB4AtDlj6GuwDWzURMD0jSg1Mq6qqovhMRkYGzH/y5AmK+O3bt8F87YZjBEMCPQQGlwOuNlih+IDUEIBphoUEUu6AexxvCNDTAeCgV8teAXZ0nOQdMH3sGGocwuI2Pj4eJTfAUv/ChQuxphUrK0hILnquAqZvTY2ABQ4jLAQGlwMcfu/DmrphPoH5FBa3sKiChQx6yMFyywFWJ+JCYNA4ACULIHFe2beilBMzF6wE82EOx6WP5BAYcAegl3DoaYDRrRnsxv+7asE03dMAzRwA8xF6yUe3cmAgHABLbygFEswhMEly6wJYVCGXgPCiEC2lDy4HwIMFmspJDQGYz2HmkBwCA+EArGkBJgirJQm1CWG1Hj6f40oDg8YBKA6BcWDtBVxFNBYfE/IoA0n9goFwALpn6d41G3QOIBT9JMsDANCtYDCjzkYBAAAAAElFTkSuQmCC"
	iconScroll    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAB30lEQVRYR2NkwA3+45BixKOHZCl8hg2YA8AWP1iRiNU3ChHzYeJUCQlshgyYA1As/vrmA9in3CICYBrGh3lfO2c9VUICOQQGlwNg3sMVEjQPAXo6ABz0BxrtwHYqqCujxDl6GoDxe5adBaubsukSRWkBlAYGhwNg+f7BzbtgH4kKC4Pp12/fooQMzOfr7z1EKSeeXvlIVkjAQ2DAHfDxQCvYB+cP7sRaAq45CykXYD43SOAF8x+chPsczL+6+jNJIQEPgYF0AMzF4MSI7hBiff7hOhNZaQKjJBxIB2ANCa2cLrA4rjiH+TxQSR4lBEKMIXWIQ/0hvGkCZ20IC4mBcABKSHj3yGBN7eg+L4kyBqsTlFBByU2wkhWpHYESIgRbRAPuAGkdfrCLBTT/gWl0n8PiGuZTmPdgIQHjv39xBxKS0JIWljYIhsBAOABcHqAHPXo+XxqqD/YRrM6A1ZIwH8JCBFZ3wNIIMSEwOBygHQop62EAVsbDouRoTRCKPC6fwxThKhdwlgODxgFYajeUOgMlGBgYGCbPXAEWstEQBNOwtIAe9wwMDGDPEwyBwegAlJISxmmJ1EXxOXrI4KoTyAkBmjsA3QJi+4C4OrMoZT96yBAsCXGkE3RzQHyyHAAA1U1oLEglgK4AAAAASUVORK5CYIIA"
	iconHelmet    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAErSURBVFhH7ZTRDcIwDES7BAPwwx9TsAB7MAbfjIHELkzADAwB2NKTnGuiEKCkSJx0H7HP51ObdvjjBdwqnBy5pZEfR7IAHM4pgeoffBuJIfhGADdY7i/O1e7k3B6vTg1AHR1z+JhhK7oF8AEMIYawtY+vLahhHgGAGumrqNUBvraghm4BXMCl4qzEOF5AI/XcjBFdqI1QEiacPIAaLdabhNTj8mhc0mcCjpAVlgwnC8BnpMbKuPwZHb6hN4I3ZhMAYsAj1HOpzln92GMLFfMIwOXRQYy1rizp9FLaQkX3AMAF/FoxwJg6OtVrgIy+isSwWwB9FRhrH1LXANp/sAoXMtg9AI8QQ/glomPuJwMAH9Ag1E0g8PonFgMf7BkARJNnjFr1VbQaNuiH4Q5NHiqUMm6nXQAAAABJRU5ErkJggg=="
)

// 16×16 pixel art stat row icons.
const (
	statIconPower     = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAB+SURBVDhPYwCC/0iYLPD/5s2bYAxiQ4RIA/83b94MxhQZ0tLSMmoIKsBlCEkGwg0B0SB+Wlraf1VVVdIN6T39/79i883/Xl5eJBsAAmDN0q455GkGORtmMwyDxCHShAGKJhANMxDEB6sgAiArBhtCbljAAIqrIEKkA6hmhv8A8o+b99NRYXUAAAAASUVORK5CYII="
	statIconIntellect = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAACXSURBVDhPY0AHLCws/wlhqFJMAJL08vL6n5aWhhOD5LEaAhJUVVUFKwLRuDBInigDApa+xkoTbQAMwzTCMEkuAOHe0///KzbfBGOSDIBhmEYYJtoAmEZ0mmgD1LJXYKXpZwAuTNAAilIiSPJAo93/BysScWKQJXgNACn6eKAVJwa5AqcBMD/iwzi9AAIgCWIwVDkQMDAAAE6lSQL2Q6wUAAAAAElFTkSuQmCC"
	statIconWisdom    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAADTSURBVDhPlZG7DcJAEETdAKIA6IBPCTSAhERAQEhOE0iQEVIEZbgJclqggkPv1rM62YuxRxr5PvPmTueqpVR4tNL7eXIzt+VhytDrsfcvZt22/8tPLktYt+1+pfqy6cB8z7v1oJIcpgRQZcCz5TSbjEVjpU99y6AseHufp8Vhkk3O4rG8JIJH3SSCKcUcQM7isUKYN9ENGZNr3FGGSliQSjDj6DZ+OjAmpF+qIt0oLNDVGVNAsA33vYXDNrVHla/HlRcNLUB5Llhu1kP92hDUgqvqC79dKEh6HYtgAAAAAElFTkSuQmCC"
	statIconAgility   = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAx0lEQVQ4T2NkoBAwUqifga4G/EdyLdxiYl3wP2Dpa7j+DdGiMDYjMQb87z39n+HwrTcMpxc0wg0papvMUGzKSFQY/Jd2zWEAaUAHhAwA+xnk9Iv33jH8OjSZwTShHmwGjP909xQMF8ADCqTRVk2Eoa8qF6wRxAZ5AwZA3kE3AOxXGABphDkbZgiyF9ANgIcyyCZkQ2DORtaMzQvggGKzywWr01cSwppAQRphABQmyF4AG0AqwAgDUg1gYGAgKiHhNZeYlIjXAAB7llkJYJkvrwAAAABJRU5ErkJgggAA"
	statIconInfluence = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAACGSURBVDhPY8AD/mPBRAGw4rS0NDj28vICY5gcSBEuANcAYqNjNIMwADG2gOWwGYLXZDRAGwNAAOx/EA3h4gS4DSDCFXgtIcoAQmoIKSDGEgxFcJoYzTDwX1VVFaw4YOlrMA3jgySJBWDNvaf//1dsvkmyZhigSDMI/FfLXkGRASCARzMDAwD4m4FwWZcjGgAAAABJRU5ErkJggg=="
)

// Environment assets for card background.
const (
	envBgSky   = "data:image/png;base64,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"
	envBgHills = "data:image/png;base64,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"
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
	Login             string
	Level             int
	NextLevel         int
	Class             string
	Title             string
	TotalXP           int
	XPInto            int
	XPFor             int
	XPBarWidth        int
	XPPct             int
	XPArcDash         int
	AccentColor       string
	ClassIcon         string
	STR               int
	INT               int
	WIS               int
	DEX               int
	CHA               int
	STRBar            int
	INTBar            int
	WISBar            int
	DEXBar            int
	CHABar            int
	StatIconPower     string
	StatIconIntellect string
	StatIconWisdom    string
	StatIconAgility   string
	StatIconInfluence string
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

const cardSVG = `<svg width="840" height="336" viewBox="0 0 700 280" xmlns="http://www.w3.org/2000/svg">
  <rect width="700" height="280" fill="#050010"/>
  <rect x="1" y="1" width="698" height="278" fill="none" stroke="#FFD700" stroke-width="2"/>
  <rect x="5" y="5" width="690" height="270" fill="none" stroke="{{.AccentColor}}" stroke-width="1"/>
  <line x1="5" y1="24" x2="5" y2="5" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="5" y1="5" x2="24" y2="5" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="676" y1="5" x2="695" y2="5" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="695" y1="5" x2="695" y2="24" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="5" y1="256" x2="5" y2="275" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="5" y1="275" x2="24" y2="275" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="676" y1="275" x2="695" y2="275" stroke="#FFD700" stroke-width="2.5"/>
  <line x1="695" y1="256" x2="695" y2="275" stroke="#FFD700" stroke-width="2.5"/>
  <image x="638" y="10" width="48" height="48" href="{{.ClassIcon}}" image-rendering="pixelated"/>
  <text x="14" y="40" font-family="Courier New,Courier,monospace" font-size="24" font-weight="bold" fill="#ffffff">{{.Login}}</text>
  <rect x="14" y="46" width="70" height="20" fill="{{.AccentColor}}"/>
  <text x="49" y="60" font-family="Courier New,Courier,monospace" font-size="11" font-weight="bold" fill="#050010" text-anchor="middle">LV.{{.Level}}</text>
  <text x="90" y="60" font-family="Courier New,Courier,monospace" font-size="13" fill="{{.AccentColor}}">{{.Class}}</text>
  <text x="14" y="74" font-family="Courier New,Courier,monospace" font-size="9" fill="#AA88DD" font-style="italic">{{.Title}}</text>
  <line x1="10" y1="80" x2="328" y2="80" stroke="#3A1A7A" stroke-width="1"/>
  <polygon points="346,75 354,80 346,85 338,80" fill="#FFD700"/>
  <line x1="362" y1="80" x2="685" y2="80" stroke="#3A1A7A" stroke-width="1"/>
  <line x1="362" y1="80" x2="362" y2="256" stroke="#3A1A7A" stroke-width="1"/>
  <image x="14" y="89" width="14" height="14" href="{{.StatIconPower}}" image-rendering="pixelated"/>
  <text x="32" y="101" font-family="Courier New,Courier,monospace" font-size="10" fill="#AA88DD">POWER</text>
  <text x="175" y="101" font-family="Courier New,Courier,monospace" font-size="13" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.STR}}</text>
  <rect x="180" y="93" width="152" height="9" fill="#180830"/>
  <rect x="180" y="93" width="{{.STRBar}}" height="9" fill="{{.AccentColor}}"/>
  <image x="14" y="117" width="14" height="14" href="{{.StatIconIntellect}}" image-rendering="pixelated"/>
  <text x="32" y="129" font-family="Courier New,Courier,monospace" font-size="10" fill="#AA88DD">INTELLECT</text>
  <text x="175" y="129" font-family="Courier New,Courier,monospace" font-size="13" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.INT}}</text>
  <rect x="180" y="121" width="152" height="9" fill="#180830"/>
  <rect x="180" y="121" width="{{.INTBar}}" height="9" fill="{{.AccentColor}}"/>
  <image x="14" y="145" width="14" height="14" href="{{.StatIconWisdom}}" image-rendering="pixelated"/>
  <text x="32" y="157" font-family="Courier New,Courier,monospace" font-size="10" fill="#AA88DD">WISDOM</text>
  <text x="175" y="157" font-family="Courier New,Courier,monospace" font-size="13" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.WIS}}</text>
  <rect x="180" y="149" width="152" height="9" fill="#180830"/>
  <rect x="180" y="149" width="{{.WISBar}}" height="9" fill="{{.AccentColor}}"/>
  <image x="14" y="173" width="14" height="14" href="{{.StatIconAgility}}" image-rendering="pixelated"/>
  <text x="32" y="185" font-family="Courier New,Courier,monospace" font-size="10" fill="#AA88DD">AGILITY</text>
  <text x="175" y="185" font-family="Courier New,Courier,monospace" font-size="13" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.DEX}}</text>
  <rect x="180" y="177" width="152" height="9" fill="#180830"/>
  <rect x="180" y="177" width="{{.DEXBar}}" height="9" fill="{{.AccentColor}}"/>
  <image x="14" y="201" width="14" height="14" href="{{.StatIconInfluence}}" image-rendering="pixelated"/>
  <text x="32" y="213" font-family="Courier New,Courier,monospace" font-size="10" fill="#AA88DD">INFLUENCE</text>
  <text x="175" y="213" font-family="Courier New,Courier,monospace" font-size="13" font-weight="bold" fill="{{.AccentColor}}" text-anchor="end">{{.CHA}}</text>
  <rect x="180" y="205" width="152" height="9" fill="#180830"/>
  <rect x="180" y="205" width="{{.CHABar}}" height="9" fill="{{.AccentColor}}"/>
  <text x="14" y="234" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD">XP {{.TotalXP}}</text>
  <text x="333" y="234" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" text-anchor="end">{{.XPInto}} / {{.XPFor}}</text>
  <rect x="14" y="238" width="320" height="7" fill="#180830"/>
  <rect x="14" y="238" width="{{.XPBarWidth}}" height="7" fill="{{.AccentColor}}"/>
  <rect x="14" y="238" width="320" height="7" fill="none" stroke="#3A1A7A" stroke-width="1"/>
  <circle cx="530" cy="155" r="62" fill="none" stroke="#3A1A7A" stroke-width="1"/>
  <circle cx="530" cy="155" r="52" fill="none" stroke="#180830" stroke-width="14"/>
  <circle cx="530" cy="155" r="52" fill="none" stroke="{{.AccentColor}}" stroke-width="14" stroke-dasharray="{{.XPArcDash}} 327" stroke-linecap="round" transform="rotate(-90 530 155)"/>
  <text x="530" y="144" font-family="Courier New,Courier,monospace" font-size="9" fill="#AA88DD" text-anchor="middle">LVL</text>
  <text x="530" y="170" font-family="Courier New,Courier,monospace" font-size="32" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.Level}}</text>
  <text x="530" y="183" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" text-anchor="middle">{{.XPPct}}%</text>
  <text x="530" y="226" font-family="Courier New,Courier,monospace" font-size="11" font-weight="bold" fill="{{.AccentColor}}" text-anchor="middle">{{.Class}}</text>
  <text x="530" y="242" font-family="Courier New,Courier,monospace" font-size="8" fill="#AA88DD" text-anchor="middle">&#8594; LV{{.NextLevel}} in {{.XPFor}} XP</text>
  <line x1="10" y1="258" x2="685" y2="258" stroke="#3A1A7A" stroke-width="1"/>
  <path transform="translate(30,259) scale(0.5)" fill="#6633CC" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
  <text x="48" y="269" font-family="Courier New,Courier,monospace" font-size="7" fill="#3A1A7A">github-rpg</text>
</svg>`

var cardTmpl = template.Must(template.New("card").Parse(cardSVG))

const (
	donutR    = 52
	donutCirc = 327 // ≈ 2π*52
)

func buildData(login string, char *stats.Character) cardData {
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
	accent := classColors[char.Class]
	if accent == "" {
		accent = classColors["Wanderer"]
	}
	icon := classIconMap[char.Class]
	if icon == "" {
		icon = iconGem
	}
	return cardData{
		Login:             login,
		Level:             char.Level,
		NextLevel:         char.Level + 1,
		Class:             char.Class,
		Title:             char.Title,
		TotalXP:           char.TotalXP,
		XPInto:            char.XPIntoLevel,
		XPFor:             char.XPForLevel,
		XPBarWidth:        xpPercent * 320 / 100,
		XPPct:             xpPercent,
		XPArcDash:         arcDash,
		AccentColor:       accent,
		ClassIcon:         icon,
		STR:               char.Strength,
		INT:               char.Intelligence,
		WIS:               char.Wisdom,
		DEX:               char.Dexterity,
		CHA:               char.Charisma,
		STRBar:            statBar(char.Strength, 152),
		INTBar:            statBar(char.Intelligence, 152),
		WISBar:            statBar(char.Wisdom, 152),
		DEXBar:            statBar(char.Dexterity, 152),
		CHABar:            statBar(char.Charisma, 152),
		StatIconPower:     statIconPower,
		StatIconIntellect: statIconIntellect,
		StatIconWisdom:    statIconWisdom,
		StatIconAgility:   statIconAgility,
		StatIconInfluence: statIconInfluence,
	}
}

// Card renders the single RPG character card as an SVG string.
// The style param is accepted for backward-compatible URLs but ignored.
func Card(login string, char *stats.Character, style string) (string, error) {
	var buf bytes.Buffer
	if err := cardTmpl.Execute(&buf, buildData(login, char)); err != nil {
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
	if err := cardTmpl.Execute(&buf, buildData(class, char)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
