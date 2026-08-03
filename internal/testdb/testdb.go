// Package testdb provides a migrated Postgres pool for store-layer tests.
// Tests skip when DATABASE_URL is unset.
package testdb

import (
	"context"
	"database/sql"
	"hash/fnv"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	migrateOnce sync.Once
	migrateErr  error
)

func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}

	migrateOnce.Do(func() { migrateErr = migrate(url) })
	if migrateErr != nil {
		t.Fatalf("migrate: %v", migrateErr)
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("pgxpool: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("postgres unreachable: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

var migrateLockID = advisoryLockID("gitrpg.testdb.migrations")

func advisoryLockID(name string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(name))
	return int64(h.Sum64())
}

func migrate(url string) error {
	dir, err := migrationsDir()
	if err != nil {
		return err
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, migrateLockID); err != nil {
		return err
	}
	defer func() {
		_, _ = conn.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, migrateLockID)
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())
	return goose.Up(db, dir)
}

func migrationsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "migrations")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// User inserts a throwaway user row and removes it when the test ends,
// cascading to every table that references it.
func User(t *testing.T, pool *pgxpool.Pool, login string) int64 {
	t.Helper()
	ctx := context.Background()
	githubID := rand.Int63n(1_000_000_000) + 8_000_000_000

	var id int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (github_id, login, access_token)
		VALUES ($1, $2, 'test-token') RETURNING id`, githubID, login).Scan(&id)
	if err != nil {
		t.Fatalf("insert test user %q: %v", login, err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id); err != nil {
			t.Logf("cleanup user %d: %v", id, err)
		}
	})
	return id
}
