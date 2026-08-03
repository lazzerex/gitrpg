package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/leaderboards"
	"github.com/lazzerex/gitrpg/internal/stats"
	"github.com/lazzerex/gitrpg/internal/users"
)

var importmapScriptRe = regexp.MustCompile(`(?s)<script type="importmap">(.*?)</script>`)

func TestImportmapHash_MatchesTemplate(t *testing.T) {
	content, err := os.ReadFile("../../web/templates/base.html")
	if err != nil {
		t.Fatalf("read base.html: %v", err)
	}
	m := importmapScriptRe.FindSubmatch(content)
	if m == nil {
		t.Fatal("no importmap script found in base.html")
	}
	digest := sha256.Sum256(m[1])
	want := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"
	if want != importmapHash {
		t.Errorf("importmapHash is stale: base.html's importmap script hashes to %s, but importmapHash = %s.\n"+
			"The CSP will silently block this script until the constant is updated to match.", want, importmapHash)
	}
}

func TestPageTemplates_ExecuteWithoutError(t *testing.T) {
	s := &Server{}
	if err := s.LoadTemplates("../../web/templates"); err != nil {
		t.Fatalf("LoadTemplates failed: %v", err)
	}

	user := &users.User{Login: "octocat"}
	char := &stats.Character{Level: 5, Class: "Guardian", Title: "The Adventurer"}

	cases := []struct {
		name string
		tmpl string
		data any
	}{
		{"index, anonymous", "index.html", indexData{BaseURL: "https://example.com"}},
		{"index, logged in", "index.html", indexData{User: user, CharClass: "Guardian", BaseURL: "https://example.com"}},
		{"profile", "profile.html", profileData{User: user, Character: char, BaseURL: "https://example.com"}},
		{"profile, needs reauth", "profile.html", profileData{User: user, Character: char, NeedsReauth: true, BaseURL: "https://example.com"}},
		{"public profile, synced", "public.html", publicProfileData{ProfileUser: user, Character: char, BaseURL: "https://example.com"}},
		{"public profile, unsynced", "public.html", publicProfileData{ProfileUser: user, BaseURL: "https://example.com"}},
		{"cards", "cards.html", baseData{BaseURL: "https://example.com"}},
		{"leaderboard", "leaderboard.html", leaderboardData{BaseURL: "https://example.com"}},
		{"game", "game.html", playData{User: user, Character: char, AccentColor: "#00add8", BaseURL: "https://example.com"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := s.templates[c.tmpl].ExecuteTemplate(&buf, "base.html", c.data); err != nil {
				t.Fatalf("execute %s: %v", c.tmpl, err)
			}
			if buf.Len() == 0 {
				t.Fatal("expected non-empty output")
			}
		})
	}
}

func TestProfile_ReauthBanner(t *testing.T) {
	s := &Server{}
	if err := s.LoadTemplates("../../web/templates"); err != nil {
		t.Fatalf("LoadTemplates failed: %v", err)
	}
	user := &users.User{Login: "octocat"}
	char := &stats.Character{Level: 5, Class: "Guardian", Title: "The Adventurer"}

	render := func(needsReauth bool) string {
		var buf bytes.Buffer
		data := profileData{User: user, Character: char, NeedsReauth: needsReauth, BaseURL: "https://example.com"}
		if err := s.templates["profile.html"].ExecuteTemplate(&buf, "base.html", data); err != nil {
			t.Fatalf("execute profile.html: %v", err)
		}
		return buf.String()
	}

	if out := render(true); !strings.Contains(out, "RECONNECT GITHUB") {
		t.Error("profile with NeedsReauth is missing the reconnect prompt")
	}
	if out := render(false); strings.Contains(out, "RECONNECT GITHUB") {
		t.Error("profile without NeedsReauth must not show the reconnect prompt")
	}
}

func TestGithubUsernameRe(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"simple", "octocat", true},
		{"with hyphen", "git-rpg", true},
		{"single char", "a", true},
		{"digits", "12345", true},
		{"max length 39", "a23456789012345678901234567890123456789"[:39], true},
		{"empty", "", false},
		{"leading hyphen", "-octocat", false},
		{"trailing hyphen", "octocat-", false},
		{"too long", "a2345678901234567890123456789012345678901", false},
		{"path traversal", "../../etc/passwd", false},
		{"sql-ish", "a' OR '1'='1", false},
		{"dot svg leftover", "octocat.svg", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := githubUsernameRe.MatchString(c.in); got != c.want {
				t.Errorf("githubUsernameRe.MatchString(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestWantsLeaderboardFragment(t *testing.T) {
	cases := []struct {
		name       string
		hxRequest  string
		hxTarget   string
		hxBoosted  string
		wantResult bool
	}{
		{"pagination click: hx-target set", "true", "leaderboard-results", "", true},
		{"boosted nav click: no hx-target", "true", "", "true", false},
		{"plain browser request", "", "", "", false},
		{"hx-request true but wrong target", "true", "some-other-id", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil)
			if c.hxRequest != "" {
				req.Header.Set("HX-Request", c.hxRequest)
			}
			if c.hxTarget != "" {
				req.Header.Set("HX-Target", c.hxTarget)
			}
			if c.hxBoosted != "" {
				req.Header.Set("HX-Boosted", c.hxBoosted)
			}
			if got := wantsLeaderboardFragment(req); got != c.wantResult {
				t.Errorf("wantsLeaderboardFragment() = %v, want %v", got, c.wantResult)
			}
		})
	}
}

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
		{"multiple entries, page 1", leaderboardData{
			Leaderboard: leaderboards.Page{
				Entries: []leaderboards.Entry{
					entry(1, "a"), entry(2, "b"), entry(3, "c"), entry(4, "d"), entry(5, "e"),
				},
				Page: 1, TotalPages: 2,
			},
		}},
		{"page 2", leaderboardData{
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
