package stats

import (
	"math"

	"github.com/lazzerex/gitrpg/internal/github"
)

// Character holds computed RPG character data derived from GitHub stats.
type Character struct {
	UserID       int64
	TotalXP      int
	Level        int
	XPIntoLevel  int
	XPForLevel   int
	Strength     int
	Intelligence int
	Wisdom       int
	Dexterity    int
	Charisma     int
	Class        string
	Title        string
}

// Calculate derives a Character from synced GitHub stats.
func Calculate(s *github.Stats) *Character {
	xp := calcXP(s)
	level, xpInto, xpFor := calcLevel(xp)
	return &Character{
		UserID:       s.UserID,
		TotalXP:      xp,
		Level:        level,
		XPIntoLevel:  xpInto,
		XPForLevel:   xpFor,
		Strength:     calcSTR(s),
		Intelligence: calcINT(s),
		Wisdom:       calcWIS(s),
		Dexterity:    calcDEX(s),
		Charisma:     calcCHA(s),
		Class:        calcClass(s),
		Title:        calcTitle(s),
	}
}

func calcXP(s *github.Stats) int {
	xp := s.CommitsCount*10 +
		s.PRsMerged*100 +
		s.IssuesClosed*25 +
		s.ReviewsCount*20 +
		s.QualifiedRepos*50
	if s.OSSReposCount >= 1 {
		xp += 500
	}
	return xp
}

// calcLevel returns the current level, XP earned above the current level
// threshold, and XP required to complete the current level.
func calcLevel(xp int) (level, xpInto, xpFor int) {
	for xpThreshold(level+1) <= xp {
		level++
	}
	var prev int
	if level > 0 {
		prev = xpThreshold(level)
	}
	next := xpThreshold(level + 1)
	return level, xp - prev, next - prev
}

// xpThreshold returns total XP required to reach level n.
// Formula: XP(N) = 100 * N^1.8, rounded to nearest integer.
func xpThreshold(n int) int {
	if n <= 0 {
		return 0
	}
	return int(math.Round(100 * math.Pow(float64(n), 1.8)))
}

func clampStat(score float64) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return int(score)
}

// STR: code output — commits and merged PRs.
// Soft caps: 3000 commits (70%), 200 PRs (30%).
func calcSTR(s *github.Stats) int {
	score := float64(s.CommitsCount)/3000*70 + float64(s.PRsMerged)/200*30
	return clampStat(score)
}

// INT: technical depth — language breadth, qualified repos, stars.
// Soft caps: 10 languages (30%), 30 qualified repos (40%), 200 stars (30%).
func calcINT(s *github.Stats) int {
	score := float64(len(s.Languages))/10*30 +
		float64(s.QualifiedRepos)/30*40 +
		float64(s.StarsReceived)/200*30
	return clampStat(score)
}

// WIS: helping others — reviews and issue closure. Floor of 5.
// Soft caps: 100 reviews (70%), 200 issues closed (30%).
func calcWIS(s *github.Stats) int {
	score := float64(s.ReviewsCount)/100*70 + float64(s.IssuesClosed)/200*30
	v := clampStat(score)
	if v < 5 {
		return 5
	}
	return v
}

// DEX: consistency — streaks and active days.
// Soft caps: 90-day current streak (30%), 365-day longest streak (40%), 90 active days (30%).
func calcDEX(s *github.Stats) int {
	score := float64(s.CurrentStreak)/90*30 +
		float64(s.LongestStreak)/365*40 +
		float64(s.ActiveDays90)/90*30
	return clampStat(score)
}

// CHA: influence — stars and followers.
// Soft caps: 500 stars (60%), 500 followers (40%).
func calcCHA(s *github.Stats) int {
	score := float64(s.StarsReceived)/500*60 + float64(s.FollowersCount)/500*40
	return clampStat(score)
}

var classMap = map[string]string{
	"Go":         "Guardian",
	"Rust":       "Berserker",
	"TypeScript": "Paladin",
	"JavaScript": "Rogue",
	"Python":     "Sage",
	"C#":         "Knight",
	"Java":       "Battlemage",
	"C++":        "Warlord",
}

func calcClass(s *github.Stats) string {
	var top string
	var max int
	for lang, count := range s.Languages {
		if count > max {
			max = count
			top = lang
		}
	}
	if class, ok := classMap[top]; ok {
		return class
	}
	return "Wanderer"
}

func calcTitle(s *github.Stats) string {
	switch {
	case s.OSSReposCount >= 5:
		return "The Open Source Hero"
	case s.QualifiedRepos >= 10 && s.StarsReceived >= 50:
		return "The Architect"
	case s.IssuesClosed >= 50:
		return "The Ticket Master"
	case s.PRsMerged >= 100:
		return "The Collaborator"
	case s.QualifiedRepos >= 5:
		return "The Maintainer"
	default:
		return "The Adventurer"
	}
}
