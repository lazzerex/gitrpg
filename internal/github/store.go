package github

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db}
}

func (s *store) upsertStats(ctx context.Context, stats *Stats) error {
	langBytes, err := json.Marshal(stats.Languages)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ctx, `
		INSERT INTO github_stats (
			user_id, commits_count, prs_merged, issues_closed, reviews_count,
			stars_received, followers_count, repos_count, qualified_repos,
			oss_repos_count, languages, longest_streak, current_streak,
			active_days_90, synced_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11::jsonb,$12,$13,$14,now(),now()
		)
		ON CONFLICT (user_id) DO UPDATE SET
			commits_count   = EXCLUDED.commits_count,
			prs_merged      = EXCLUDED.prs_merged,
			issues_closed   = EXCLUDED.issues_closed,
			reviews_count   = EXCLUDED.reviews_count,
			stars_received  = EXCLUDED.stars_received,
			followers_count = EXCLUDED.followers_count,
			repos_count     = EXCLUDED.repos_count,
			qualified_repos = EXCLUDED.qualified_repos,
			oss_repos_count = EXCLUDED.oss_repos_count,
			languages       = EXCLUDED.languages,
			longest_streak  = EXCLUDED.longest_streak,
			current_streak  = EXCLUDED.current_streak,
			active_days_90  = EXCLUDED.active_days_90,
			synced_at       = EXCLUDED.synced_at,
			updated_at      = EXCLUDED.updated_at
	`,
		stats.UserID, stats.CommitsCount, stats.PRsMerged, stats.IssuesClosed,
		stats.ReviewsCount, stats.StarsReceived, stats.FollowersCount,
		stats.ReposCount, stats.QualifiedRepos, stats.OSSReposCount,
		langBytes, stats.LongestStreak, stats.CurrentStreak, stats.ActiveDays90,
	)
	return err
}

func (s *store) startSync(ctx context.Context, userID int64) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx,
		`INSERT INTO github_syncs (user_id) VALUES ($1) RETURNING id`,
		userID,
	).Scan(&id)
	return id, err
}

func (s *store) completeSync(ctx context.Context, id int64, pointsUsed int, syncErr error) error {
	status := "success"
	errMsg := ""
	if syncErr != nil {
		status = "failed"
		errMsg = syncErr.Error()
	}
	_, err := s.db.Exec(ctx, `
		UPDATE github_syncs
		SET completed_at=$1, status=$2, points_used=$3, error=NULLIF($4,'')
		WHERE id=$5
	`, time.Now(), status, pointsUsed, errMsg, id)
	return err
}

func (s *store) getStats(ctx context.Context, userID int64) (*Stats, error) {
	var stats Stats
	var langStr string

	err := s.db.QueryRow(ctx, `
		SELECT user_id, commits_count, prs_merged, issues_closed, reviews_count,
		       stars_received, followers_count, repos_count, qualified_repos,
		       oss_repos_count, languages::text, longest_streak, current_streak, active_days_90
		FROM github_stats WHERE user_id=$1
	`, userID).Scan(
		&stats.UserID, &stats.CommitsCount, &stats.PRsMerged, &stats.IssuesClosed,
		&stats.ReviewsCount, &stats.StarsReceived, &stats.FollowersCount,
		&stats.ReposCount, &stats.QualifiedRepos, &stats.OSSReposCount,
		&langStr, &stats.LongestStreak, &stats.CurrentStreak, &stats.ActiveDays90,
	)
	if err != nil {
		return nil, err
	}

	stats.Languages = make(map[string]int)
	_ = json.Unmarshal([]byte(langStr), &stats.Languages)

	return &stats, nil
}
