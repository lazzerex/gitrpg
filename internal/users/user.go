package users

import "time"

type User struct {
	ID          int64
	GitHubID    int64
	Login       string
	Name        string
	AvatarURL   string
	Email       string
	AccessToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type contextKey struct{}

var ContextKey = contextKey{}
