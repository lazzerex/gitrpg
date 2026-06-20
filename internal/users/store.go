package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Upsert(ctx context.Context, u *User) (*User, error) {
	const q = `
		INSERT INTO users (github_id, login, name, avatar_url, email, access_token, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (github_id) DO UPDATE SET
			login        = EXCLUDED.login,
			name         = EXCLUDED.name,
			avatar_url   = EXCLUDED.avatar_url,
			email        = EXCLUDED.email,
			access_token = EXCLUDED.access_token,
			updated_at   = now()
		RETURNING id, github_id, login, name, avatar_url, email, access_token, created_at, updated_at
	`
	row := s.db.QueryRow(ctx, q, u.GitHubID, u.Login, u.Name, u.AvatarURL, u.Email, u.AccessToken)
	return scanUser(row)
}

func (s *Store) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `
		SELECT id, github_id, login, name, avatar_url, email, access_token, created_at, updated_at
		FROM users WHERE id = $1
	`
	row := s.db.QueryRow(ctx, q, id)
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func scanUser(row pgx.Row) (*User, error) {
	var u User
	err := row.Scan(
		&u.ID, &u.GitHubID, &u.Login, &u.Name,
		&u.AvatarURL, &u.Email, &u.AccessToken,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
