package github

import (
	"sort"
	"strings"
	"time"
)

// Stats holds processed RPG-relevant statistics derived from GitHub data.
type Stats struct {
	UserID         int64
	CommitsCount   int
	PRsMerged      int
	IssuesClosed   int
	ReviewsCount   int
	StarsReceived  int
	FollowersCount int
	ReposCount     int
	QualifiedRepos int
	OSSReposCount  int
	Languages      map[string]int // language name → commit count
	LongestStreak  int
	CurrentStreak  int
	ActiveDays90   int
}

func process(userID int64, raw *RawStats) *Stats {
	s := &Stats{
		UserID:         userID,
		CommitsCount:   raw.Commits,
		PRsMerged:      raw.PRsMerged,
		IssuesClosed:   raw.IssuesClosed,
		ReviewsCount:   raw.Reviews,
		FollowersCount: raw.Followers,
		Languages:      make(map[string]int),
	}

	prefix := raw.Login + "/"

	// Merge repos from both sources; deduplicate by nameWithOwner.
	type repoData struct {
		IsFork         bool
		StargazerCount int
		ForkCount      int
		CommitCount    int
	}
	merged := make(map[string]repoData)

	for _, r := range raw.RepoContribs {
		if r.Language != "" {
			s.Languages[r.Language] += r.CommitCount
		}
		d := merged[r.NameWithOwner]
		d.IsFork = r.IsFork
		d.StargazerCount = r.StargazerCount
		d.ForkCount = r.ForkCount
		d.CommitCount += r.CommitCount
		merged[r.NameWithOwner] = d
	}

	// AllRepos provides authoritative star/fork counts for all owned repos.
	for _, r := range raw.AllRepos {
		d := merged[r.NameWithOwner]
		d.IsFork = r.IsFork
		d.StargazerCount = r.StargazerCount
		d.ForkCount = r.ForkCount
		merged[r.NameWithOwner] = d
	}

	ossSet := make(map[string]struct{})
	qualifiedSet := make(map[string]struct{})

	for name, d := range merged {
		isOwned := strings.HasPrefix(name, prefix)

		if !isOwned && d.CommitCount > 0 {
			ossSet[name] = struct{}{}
			continue
		}

		if isOwned && !d.IsFork {
			s.ReposCount++
			s.StarsReceived += d.StargazerCount
			if d.StargazerCount >= 1 || d.ForkCount >= 1 || d.CommitCount >= 5 {
				qualifiedSet[name] = struct{}{}
			}
		}
	}

	s.QualifiedRepos = len(qualifiedSet)
	s.OSSReposCount = len(ossSet)
	s.LongestStreak, s.CurrentStreak, s.ActiveDays90 = calculateStreaks(raw.Calendar)

	return s
}

// Year windows are day-aligned but contributionCalendar returns whole weeks,
// so boundary days arrive once per window.
func dedupeDays(days []CalendarDay) []CalendarDay {
	best := make(map[string]int, len(days))
	for _, d := range days {
		if c, ok := best[d.Date]; !ok || d.Count > c {
			best[d.Date] = d.Count
		}
	}
	out := make([]CalendarDay, 0, len(best))
	for date, count := range best {
		out = append(out, CalendarDay{Date: date, Count: count})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	return out
}

func calculateStreaks(days []CalendarDay) (longest, current, activeDays90 int) {
	days = dedupeDays(days)
	now := time.Now()
	cutoff90 := now.AddDate(0, 0, -90)
	today := now.Format("2006-01-02")

	streak := 0
	for _, d := range days {
		t, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		if d.Count > 0 {
			streak++
			if streak > longest {
				longest = streak
			}
			if !t.Before(cutoff90) {
				activeDays90++
			}
		} else {
			streak = 0
		}
	}

	// Walk backwards for current streak; skip today if no contributions yet.
	for i := len(days) - 1; i >= 0; i-- {
		d := days[i]
		if d.Date == today && d.Count == 0 {
			continue
		}
		if d.Count > 0 {
			current++
		} else {
			break
		}
	}

	return
}
