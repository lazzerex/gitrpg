package achievements

import "github.com/lazzerex/gitrpg/internal/github"

type Rarity string

const (
	Common    Rarity = "common"
	Rare      Rarity = "rare"
	Legendary Rarity = "legendary"
)

type Achievement struct {
	Slug        string
	Name        string
	Description string
	Icon        string
	Rarity      Rarity
	check       func(s *github.Stats) bool
}

// RarityColor returns the Tailwind text color class for this achievement's rarity.
func (a Achievement) RarityColor() string {
	switch a.Rarity {
	case Legendary:
		return "text-yellow-400"
	case Rare:
		return "text-blue-400"
	default:
		return "text-gray-300"
	}
}

// AchClass returns CSS class for the achievement border/shadow based on rarity and earned state.
func (a Achievement) AchClass() string {
	switch a.Rarity {
	case Legendary:
		return "ach-legendary"
	case Rare:
		return "ach-rare"
	default:
		return "ach-common"
	}
}

// UserAchievement wraps an Achievement with the user's earned status.
type UserAchievement struct {
	Achievement
	Earned bool
}

// All is the canonical list of achievements, in display order.
var All = []Achievement{
	// Common
	{
		Slug: "first-commit", Name: "First Commit", Icon: "git-commit-horizontal", Rarity: Common,
		Description: "Made your first commit.",
		check:       func(s *github.Stats) bool { return s.CommitsCount >= 1 },
	},
	{
		Slug: "first-pr", Name: "First Pull Request", Icon: "git-pull-request", Rarity: Common,
		Description: "Merged your first pull request.",
		check:       func(s *github.Stats) bool { return s.PRsMerged >= 1 },
	},
	{
		Slug: "first-repo", Name: "First Repository", Icon: "folder-git-2", Rarity: Common,
		Description: "Created your first qualified repository.",
		check:       func(s *github.Stats) bool { return s.QualifiedRepos >= 1 },
	},
	// Rare
	{
		Slug: "commits-1k", Name: "Thousand Commits", Icon: "flame", Rarity: Rare,
		Description: "Reached 1,000 total commits.",
		check:       func(s *github.Stats) bool { return s.CommitsCount >= 1000 },
	},
	{
		Slug: "prs-100", Name: "Century of PRs", Icon: "git-merge", Rarity: Rare,
		Description: "Merged 100 pull requests.",
		check:       func(s *github.Stats) bool { return s.PRsMerged >= 100 },
	},
	{
		Slug: "reviews-10", Name: "Code Reviewer", Icon: "eye", Rarity: Rare,
		Description: "Submitted 10 pull request reviews.",
		check:       func(s *github.Stats) bool { return s.ReviewsCount >= 10 },
	},
	{
		Slug: "streak-365", Name: "Year-Long Streak", Icon: "calendar-check", Rarity: Rare,
		Description: "Maintained a 365-day contribution streak.",
		check:       func(s *github.Stats) bool { return s.LongestStreak >= 365 },
	},
	// Legendary
	{
		Slug: "oss-hero", Name: "Open Source Hero", Icon: "globe", Rarity: Legendary,
		Description: "Contributed to 5 or more open source repositories.",
		check:       func(s *github.Stats) bool { return s.OSSReposCount >= 5 },
	},
	{
		Slug: "stars-10k", Name: "Star Collector", Icon: "star", Rarity: Legendary,
		Description: "Received 10,000 stars across repositories.",
		check:       func(s *github.Stats) bool { return s.StarsReceived >= 10000 },
	},
}

// Evaluate returns slugs of all achievements earned by the given stats.
func Evaluate(s *github.Stats) []string {
	var earned []string
	for _, a := range All {
		if a.check(s) {
			earned = append(earned, a.Slug)
		}
	}
	return earned
}

// BuildUserAchievements merges All with the user's earned slug set.
func BuildUserAchievements(earnedSlugs []string) []UserAchievement {
	earned := make(map[string]bool, len(earnedSlugs))
	for _, s := range earnedSlugs {
		earned[s] = true
	}
	out := make([]UserAchievement, len(All))
	for i, a := range All {
		out[i] = UserAchievement{Achievement: a, Earned: earned[a.Slug]}
	}
	return out
}
