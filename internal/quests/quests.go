package quests

import (
	"time"

	"github.com/lazzerex/gitrpg/internal/github"
)

// Quest is a weekly challenge defined against a cumulative GitHub stat counter.
// Progress is the counter's delta since the quest was assigned.
type Quest struct {
	Slug   string
	Name   string
	Target int
	XP     int
	metric func(*github.Stats) int
}

var catalog = []Quest{
	{Slug: "weekly-prs", Name: "Merge 2 pull requests", Target: 2, XP: 150,
		metric: func(s *github.Stats) int { return s.PRsMerged }},
	{Slug: "weekly-reviews", Name: "Submit 3 reviews", Target: 3, XP: 100,
		metric: func(s *github.Stats) int { return s.ReviewsCount }},
	{Slug: "weekly-issues", Name: "Close 2 issues", Target: 2, XP: 75,
		metric: func(s *github.Stats) int { return s.IssuesClosed }},
	{Slug: "weekly-commits", Name: "Land 20 commits", Target: 20, XP: 100,
		metric: func(s *github.Stats) int { return s.CommitsCount }},
}

// weekStart returns the Monday 00:00 UTC of the week containing t.
func weekStart(t time.Time) time.Time {
	t = t.UTC()
	back := (int(t.Weekday()) + 6) % 7
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -back)
}

// progressOf clamps the counter delta to [0, target]. A negative delta can
// occur when a re-sync corrects historical counts downward.
func progressOf(current, baseline, target int) int {
	p := current - baseline
	if p < 0 {
		p = 0
	}
	if p > target {
		p = target
	}
	return p
}
