package server

import (
	"bytes"
	"testing"

	"github.com/lazzerex/gitrpg/internal/leaderboards"
	"github.com/lazzerex/gitrpg/internal/users"
)

func TestLoadTemplates(t *testing.T) {
	s := &Server{}
	if err := s.LoadTemplates("../../web/templates"); err != nil {
		t.Fatalf("LoadTemplates failed: %v", err)
	}
	for _, name := range []string{"index.html", "profile.html", "public.html", "cards.html", "leaderboard.html", "char-panel"} {
		if _, ok := s.templates[name]; !ok {
			t.Errorf("template %q not registered", name)
		}
	}
}

func TestLeaderboardResultsFragmentExecutes(t *testing.T) {
	s := &Server{}
	if err := s.LoadTemplates("../../web/templates"); err != nil {
		t.Fatalf("LoadTemplates failed: %v", err)
	}

	entry := func(rank int, login string) leaderboards.Entry {
		return leaderboards.Entry{Rank: rank, Login: login, Class: "Guardian", Level: 5, TotalXP: 1000}
	}

	cases := []struct {
		name string
		data leaderboardData
	}{
		{"single entry, page 1", leaderboardData{
			Leaderboard: leaderboards.Page{Entries: []leaderboards.Entry{entry(1, "octocat")}, Page: 1, TotalPages: 1},
		}},
		{"full podium plus overflow list, page 1", leaderboardData{
			Leaderboard: leaderboards.Page{
				Entries: []leaderboards.Entry{
					entry(1, "a"), entry(2, "b"), entry(3, "c"), entry(4, "d"), entry(5, "e"),
				},
				Page: 1, TotalPages: 2,
			},
		}},
		{"page 2, no podium", leaderboardData{
			Leaderboard: leaderboards.Page{Entries: []leaderboards.Entry{entry(51, "z")}, Page: 2, TotalPages: 2},
		}},
		{"empty leaderboard", leaderboardData{
			Leaderboard: leaderboards.Page{Entries: nil, Page: 1, TotalPages: 1},
		}},
		{"viewer present, gets YOU highlight", leaderboardData{
			User:        &users.User{Login: "octocat"},
			Leaderboard: leaderboards.Page{Entries: []leaderboards.Entry{entry(1, "octocat"), entry(2, "b"), entry(3, "c"), entry(4, "octocat")}, Page: 1, TotalPages: 1},
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := s.templates["leaderboard.html"].ExecuteTemplate(&buf, "leaderboard-results", c.data); err != nil {
				t.Fatalf("execute leaderboard-results: %v", err)
			}
			if buf.Len() == 0 {
				t.Fatal("expected non-empty fragment output")
			}
			if c.name == "viewer present, gets YOU highlight" {
				out := buf.String()
				if !bytes.Contains(buf.Bytes(), []byte("lb-row-you")) {
					t.Errorf("expected lb-row-you class in output, got:\n%s", out)
				}
			}
		})
	}
}
