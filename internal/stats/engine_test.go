package stats

import (
	"testing"

	"github.com/lazzerex/gitrpg/internal/github"
)

func TestCalcXP(t *testing.T) {
	cases := []struct {
		name string
		s    *github.Stats
		want int
	}{
		{"empty", &github.Stats{}, 0},
		{"commits only", &github.Stats{CommitsCount: 5}, 50},
		{"oss bonus", &github.Stats{OSSReposCount: 1}, 500},
		{"mixed", &github.Stats{
			CommitsCount: 10, PRsMerged: 2, IssuesClosed: 4,
			ReviewsCount: 3, QualifiedRepos: 1,
		}, 10*10 + 2*100 + 4*25 + 3*20 + 1*50},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := calcXP(c.s); got != c.want {
				t.Errorf("calcXP() = %d, want %d", got, c.want)
			}
		})
	}
}

func TestCalcLevel(t *testing.T) {
	cases := []struct {
		name      string
		xp        int
		wantLevel int
	}{
		{"zero xp is level 0", 0, 0},
		{"just under level 1", xpThreshold(1) - 1, 0},
		{"exactly level 1", xpThreshold(1), 1},
		{"level 5", xpThreshold(5), 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			level, xpInto, xpFor := calcLevel(c.xp)
			if level != c.wantLevel {
				t.Errorf("calcLevel(%d) level = %d, want %d", c.xp, level, c.wantLevel)
			}
			if xpInto < 0 || xpInto >= xpFor && xpFor != 0 {
				t.Errorf("calcLevel(%d) xpInto=%d xpFor=%d out of range", c.xp, xpInto, xpFor)
			}
		})
	}
}

func TestClampStat(t *testing.T) {
	cases := []struct {
		score float64
		want  int
	}{
		{-10, 0},
		{0, 0},
		{50.7, 50},
		{100, 100},
		{200, 100},
	}
	for _, c := range cases {
		if got := clampStat(c.score); got != c.want {
			t.Errorf("clampStat(%v) = %d, want %d", c.score, got, c.want)
		}
	}
}

func TestCalcWIS_Floor(t *testing.T) {
	got := calcWIS(&github.Stats{})
	if got != 5 {
		t.Errorf("calcWIS(empty) = %d, want floor 5", got)
	}
}

func TestCalcClass(t *testing.T) {
	cases := []struct {
		name string
		s    *github.Stats
		want string
	}{
		{"no languages", &github.Stats{}, "Wanderer"},
		{"unknown language", &github.Stats{Languages: map[string]int{"COBOL": 100}}, "Wanderer"},
		{"go dominant", &github.Stats{Languages: map[string]int{"Go": 100, "Python": 10}}, "Guardian"},
		{"python dominant", &github.Stats{Languages: map[string]int{"Go": 10, "Python": 100}}, "Sage"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := calcClass(c.s); got != c.want {
				t.Errorf("calcClass() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestCalcTitle(t *testing.T) {
	cases := []struct {
		name string
		s    *github.Stats
		want string
	}{
		{"default", &github.Stats{}, "The Adventurer"},
		{"maintainer", &github.Stats{QualifiedRepos: 5}, "The Maintainer"},
		{"collaborator", &github.Stats{PRsMerged: 100}, "The Collaborator"},
		{"ticket master", &github.Stats{IssuesClosed: 50}, "The Ticket Master"},
		{"architect", &github.Stats{QualifiedRepos: 10, StarsReceived: 50}, "The Architect"},
		{"oss hero takes priority", &github.Stats{OSSReposCount: 5, QualifiedRepos: 10, StarsReceived: 50}, "The Open Source Hero"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := calcTitle(c.s); got != c.want {
				t.Errorf("calcTitle() = %q, want %q", got, c.want)
			}
		})
	}
}
