package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/crypto"
)

var ErrNotFound = errors.New("user not found")

type Store struct {
	db  *pgxpool.Pool
	key []byte // nil = no encryption (dev)
}

func NewStore(db *pgxpool.Pool, key []byte) *Store {
	return &Store{db: db, key: key}
}

func (s *Store) Upsert(ctx context.Context, u *User) (*User, error) {
	storedToken, err := s.encryptToken(u.AccessToken)
	if err != nil {
		return nil, err
	}

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
	row := s.db.QueryRow(ctx, q, u.GitHubID, u.Login, u.Name, u.AvatarURL, u.Email, storedToken)
	return s.scanUser(row)
}

func (s *Store) GetByLogin(ctx context.Context, login string) (*User, error) {
	const q = `
		SELECT id, github_id, login, name, avatar_url, email, access_token, created_at, updated_at
		FROM users WHERE login = $1
	`
	row := s.db.QueryRow(ctx, q, login)
	u, err := s.scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (s *Store) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `
		SELECT id, github_id, login, name, avatar_url, email, access_token, created_at, updated_at
		FROM users WHERE id = $1
	`
	row := s.db.QueryRow(ctx, q, id)
	u, err := s.scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (s *Store) ListAll(ctx context.Context) ([]*User, error) {
	const q = `
		SELECT id, github_id, login, name, avatar_url, email, access_token, created_at, updated_at
		FROM users ORDER BY id
	`
	rows, err := s.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*User
	for rows.Next() {
		var u User
		var storedToken string
		if err := rows.Scan(
			&u.ID, &u.GitHubID, &u.Login, &u.Name,
			&u.AvatarURL, &u.Email, &storedToken,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.AccessToken = s.decryptToken(storedToken)
		out = append(out, &u)
	}
	return out, rows.Err()
}

func (s *Store) scanUser(row pgx.Row) (*User, error) {
	var u User
	var storedToken string
	err := row.Scan(
		&u.ID, &u.GitHubID, &u.Login, &u.Name,
		&u.AvatarURL, &u.Email, &storedToken,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.AccessToken = s.decryptToken(storedToken)
	return &u, nil
}

func (s *Store) encryptToken(plaintext string) (string, error) {
	if len(s.key) == 0 {
		return plaintext, nil
	}
	return crypto.Seal([]byte(plaintext), s.key)
}

// decryptToken decrypts a stored token. If decryption fails (e.g. legacy plaintext),
// returns the raw value so the user can still authenticate until next login re-encrypts.
func (s *Store) decryptToken(stored string) string {
	if len(s.key) == 0 {
		return stored
	}
	plaintext, err := crypto.Open(stored, s.key)
	if err != nil {
		// legacy plaintext token — will be re-encrypted on next login
		return stored
	}
	return string(plaintext)
}
