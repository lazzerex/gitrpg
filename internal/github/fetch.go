package github

import (
	"context"
	"fmt"
	"log/slog"
)

// contributionsCollection covers only the most recent year.
// TODO: loop over prior years for all-time commit accuracy.
const statsQuery = `
query($login: String!) {
  rateLimit { cost remaining }
  user(login: $login) {
    followers    { totalCount }
    pullRequests(states: [MERGED]) { totalCount }
    issues      (states: [CLOSED]) { totalCount }
    contributionsCollection {
      totalCommitContributions
      totalPullRequestReviewContributions
      contributionCalendar {
        weeks {
          contributionDays { contributionCount date }
        }
      }
      commitContributionsByRepository(maxRepositories: 100) {
        repository {
          nameWithOwner
          isFork
          stargazerCount
          forkCount
          primaryLanguage { name }
        }
        contributions { totalCount }
      }
    }
  }
}`

const reposQuery = `
query($login: String!, $after: String) {
  rateLimit { cost remaining }
  user(login: $login) {
    repositories(first: 100, after: $after, ownerAffiliations: [OWNER]) {
      pageInfo { hasNextPage endCursor }
      nodes {
        nameWithOwner
        isFork
        stargazerCount
        forkCount
        primaryLanguage { name }
      }
    }
  }
}`

type statsResult struct {
	RateLimit struct {
		Cost      int `json:"cost"`
		Remaining int `json:"remaining"`
	} `json:"rateLimit"`
	User *struct {
		Followers    struct{ TotalCount int } `json:"followers"`
		PullRequests struct{ TotalCount int } `json:"pullRequests"`
		Issues       struct{ TotalCount int } `json:"issues"`
		ContributionsCollection struct {
			TotalCommitContributions            int `json:"totalCommitContributions"`
			TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions"`
			ContributionCalendar                struct {
				Weeks []struct {
					ContributionDays []struct {
						ContributionCount int    `json:"contributionCount"`
						Date              string `json:"date"`
					} `json:"contributionDays"`
				} `json:"weeks"`
			} `json:"contributionCalendar"`
			CommitContributionsByRepository []struct {
				Repository struct {
					NameWithOwner   string  `json:"nameWithOwner"`
					IsFork          bool    `json:"isFork"`
					StargazerCount  int     `json:"stargazerCount"`
					ForkCount       int     `json:"forkCount"`
					PrimaryLanguage *struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
				} `json:"repository"`
				Contributions struct{ TotalCount int } `json:"contributions"`
			} `json:"commitContributionsByRepository"`
		} `json:"contributionsCollection"`
	} `json:"user"`
}

type reposResult struct {
	RateLimit struct {
		Cost int `json:"cost"`
	} `json:"rateLimit"`
	User *struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
			Nodes []struct {
				NameWithOwner   string `json:"nameWithOwner"`
				IsFork          bool   `json:"isFork"`
				StargazerCount  int    `json:"stargazerCount"`
				ForkCount       int    `json:"forkCount"`
				PrimaryLanguage *struct {
					Name string `json:"name"`
				} `json:"primaryLanguage"`
			} `json:"nodes"`
		} `json:"repositories"`
	} `json:"user"`
}

// RawRepo is a repository record from GitHub.
type RawRepo struct {
	NameWithOwner  string
	IsFork         bool
	StargazerCount int
	ForkCount      int
	Language       string
	CommitCount    int
}

// CalendarDay is a single day in the contribution calendar.
type CalendarDay struct {
	Date  string
	Count int
}

// RawStats holds unprocessed data fetched from GitHub.
type RawStats struct {
	Login        string
	Followers    int
	PRsMerged    int
	IssuesClosed int
	Commits      int
	Reviews      int
	RepoContribs []RawRepo
	AllRepos     []RawRepo
	Calendar     []CalendarDay
	PointsUsed   int
}

func fetch(ctx context.Context, token, login string, logger *slog.Logger) (*RawStats, error) {
	c := newClient(token, logger)

	var sr statsResult
	if err := c.query(ctx, statsQuery, map[string]any{"login": login}, &sr); err != nil {
		return nil, err
	}
	if sr.User == nil {
		return nil, fmt.Errorf("github user not found: %s", login)
	}

	raw := &RawStats{
		Login:        login,
		Followers:    sr.User.Followers.TotalCount,
		PRsMerged:    sr.User.PullRequests.TotalCount,
		IssuesClosed: sr.User.Issues.TotalCount,
		Commits:      sr.User.ContributionsCollection.TotalCommitContributions,
		Reviews:      sr.User.ContributionsCollection.TotalPullRequestReviewContributions,
		PointsUsed:   sr.RateLimit.Cost,
	}

	cc := sr.User.ContributionsCollection
	for _, w := range cc.ContributionCalendar.Weeks {
		for _, d := range w.ContributionDays {
			raw.Calendar = append(raw.Calendar, CalendarDay{Date: d.Date, Count: d.ContributionCount})
		}
	}

	for _, rc := range cc.CommitContributionsByRepository {
		r := RawRepo{
			NameWithOwner:  rc.Repository.NameWithOwner,
			IsFork:         rc.Repository.IsFork,
			StargazerCount: rc.Repository.StargazerCount,
			ForkCount:      rc.Repository.ForkCount,
			CommitCount:    rc.Contributions.TotalCount,
		}
		if rc.Repository.PrimaryLanguage != nil {
			r.Language = rc.Repository.PrimaryLanguage.Name
		}
		raw.RepoContribs = append(raw.RepoContribs, r)
	}

	// Paginate owned repos for complete star count (up to 2000 repos).
	var after string
	for page := 0; page < 20; page++ {
		vars := map[string]any{"login": login}
		if after != "" {
			vars["after"] = after
		}
		var rr reposResult
		if err := c.query(ctx, reposQuery, vars, &rr); err != nil {
			return nil, err
		}
		if rr.User == nil {
			break
		}
		raw.PointsUsed += rr.RateLimit.Cost
		for _, node := range rr.User.Repositories.Nodes {
			r := RawRepo{
				NameWithOwner:  node.NameWithOwner,
				IsFork:         node.IsFork,
				StargazerCount: node.StargazerCount,
				ForkCount:      node.ForkCount,
			}
			if node.PrimaryLanguage != nil {
				r.Language = node.PrimaryLanguage.Name
			}
			raw.AllRepos = append(raw.AllRepos, r)
		}
		if !rr.User.Repositories.PageInfo.HasNextPage {
			break
		}
		after = rr.User.Repositories.PageInfo.EndCursor
	}

	return raw, nil
}
