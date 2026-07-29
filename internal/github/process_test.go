package github

import (
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	raw := &RawStats{
		Login:        "alice",
		Commits:      50,
		PRsMerged:    3,
		IssuesClosed: 2,
		Reviews:      1,
		Followers:    10,
		RepoContribs: []RawRepo{
			{NameWithOwner: "alice/own-repo", Language: "Go", CommitCount: 10},
			{NameWithOwner: "someoneelse/oss-repo", Language: "Python", CommitCount: 5},
		},
		AllRepos: []RawRepo{
			{NameWithOwner: "alice/own-repo", StargazerCount: 2, ForkCount: 0},
			{NameWithOwner: "alice/forked-repo", IsFork: true},
		},
	}

	s := process(1, raw)

	if s.UserID != 1 {
		t.Errorf("UserID = %d, want 1", s.UserID)
	}
	if s.CommitsCount != 50 {
		t.Errorf("CommitsCount = %d, want 50", s.CommitsCount)
	}
	if s.Languages["Go"] != 10 || s.Languages["Python"] != 5 {
		t.Errorf("Languages = %+v, want Go:10 Python:5", s.Languages)
	}
	if s.ReposCount != 1 {
		t.Errorf("ReposCount = %d, want 1 (fork excluded)", s.ReposCount)
	}
	if s.StarsReceived != 2 {
		t.Errorf("StarsReceived = %d, want 2", s.StarsReceived)
	}
	if s.QualifiedRepos != 1 {
		t.Errorf("QualifiedRepos = %d, want 1", s.QualifiedRepos)
	}
	if s.OSSReposCount != 1 {
		t.Errorf("OSSReposCount = %d, want 1 (external repo with commits)", s.OSSReposCount)
	}
}

func TestCalculateStreaks_Longest(t *testing.T) {
	days := []CalendarDay{
		{Date: "2020-01-01", Count: 1},
		{Date: "2020-01-02", Count: 1},
		{Date: "2020-01-03", Count: 1},
		{Date: "2020-01-04", Count: 0},
		{Date: "2020-01-05", Count: 1},
	}
	longest, _, _ := calculateStreaks(days)
	if longest != 3 {
		t.Errorf("longest = %d, want 3", longest)
	}
}

func TestCalculateStreaks_Current(t *testing.T) {
	now := time.Now()
	days := []CalendarDay{
		{Date: now.AddDate(0, 0, -2).Format("2006-01-02"), Count: 1},
		{Date: now.AddDate(0, 0, -1).Format("2006-01-02"), Count: 1},
		{Date: now.Format("2006-01-02"), Count: 1},
	}
	_, current, _ := calculateStreaks(days)
	if current != 3 {
		t.Errorf("current = %d, want 3", current)
	}
}

func TestCalculateStreaks_CurrentSkipsEmptyToday(t *testing.T) {
	now := time.Now()
	days := []CalendarDay{
		{Date: now.AddDate(0, 0, -2).Format("2006-01-02"), Count: 1},
		{Date: now.AddDate(0, 0, -1).Format("2006-01-02"), Count: 1},
		{Date: now.Format("2006-01-02"), Count: 0},
	}
	_, current, _ := calculateStreaks(days)
	if current != 2 {
		t.Errorf("current = %d, want 2 (today not yet contributed)", current)
	}
}

func TestCalculateStreaks_ActiveDays90(t *testing.T) {
	now := time.Now()
	days := []CalendarDay{
		{Date: now.AddDate(0, 0, -100).Format("2006-01-02"), Count: 1},
		{Date: now.AddDate(0, 0, -10).Format("2006-01-02"), Count: 1},
	}
	_, _, active := calculateStreaks(days)
	if active != 1 {
		t.Errorf("activeDays90 = %d, want 1 (100-day-old entry outside window)", active)
	}
}
