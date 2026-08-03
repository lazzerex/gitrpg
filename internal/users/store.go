package users

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lazzerex/gitrpg/internal/crypto"
)

var ErrNotFound = errors.New("user not found")

type Store struct {
	db     *pgxpool.Pool
	key    []byte // nil = no encryption (dev)
	logger *slog.Logger
}

func NewStore(db *pgxpool.Pool, key []byte, logger *slog.Logger) *Store {
	return &Store{db: db, key: key, logger: logger}
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
		RETURNING id, github_id, login, COALESCE(name,''), COALESCE(avatar_url,''), COALESCE(email,''), access_token, created_at, updated_at
	`
	row := s.db.QueryRow(ctx, q, u.GitHubID, u.Login, u.Name, u.AvatarURL, u.Email, storedToken)
	return s.scanUser(row)
}

// Logins are renameable and recyclable, so two rows can hold one login until
// the renamed account signs in again; most recently updated is the owner.
func (s *Store) GetByLogin(ctx context.Context, login string) (*User, error) {
	const q = `
		SELECT id, github_id, login, COALESCE(name,''), COALESCE(avatar_url,''), COALESCE(email,''), access_token, created_at, updated_at
		FROM users WHERE login = $1
		ORDER BY updated_at DESC, id DESC
		LIMIT 1
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
		SELECT id, github_id, login, COALESCE(name,''), COALESCE(avatar_url,''), COALESCE(email,''), access_token, created_at, updated_at
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
		SELECT id, github_id, login, COALESCE(name,''), COALESCE(avatar_url,''), COALESCE(email,''), access_token, created_at, updated_at
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
		if !isPlaintextToken(stored) {
			s.logger.Warn("stored access token failed to decrypt, sending it as-is will fail; TOKEN_ENCRYPTION_KEY was likely rotated and the user must sign in again")
		}
		return stored
	}
	return string(plaintext)
}

var tokenPrefixes = []string{"gho_", "ghu_", "ghp_", "ghs_", "github_pat_"}

func isPlaintextToken(s string) bool {
	for _, p := range tokenPrefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
