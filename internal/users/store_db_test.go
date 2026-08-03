package users

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/lazzerex/gitrpg/internal/testdb"
)

func dbStore(t *testing.T, key []byte) *Store {
	t.Helper()
	return NewStore(testdb.Pool(t), key, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestStore_UpsertAndGet(t *testing.T) {
	s := dbStore(t, nil)
	ctx := context.Background()
	login := "store-upsert-" + t.Name()

	u, err := s.Upsert(ctx, &User{GitHubID: 8100000001, Login: login, Name: "First", AccessToken: "gho_first"})
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	t.Cleanup(func() { _, _ = s.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, u.ID) })
	if u.ID == 0 {
		t.Fatal("Upsert returned zero id")
	}

	got, err := s.GetByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Login != login || got.Name != "First" || got.AccessToken != "gho_first" {
		t.Errorf("GetByID = %+v, want login %q name First token gho_first", got, login)
	}

	renamed := login + "-renamed"
	u2, err := s.Upsert(ctx, &User{GitHubID: 8100000001, Login: renamed, Name: "Second", AccessToken: "gho_second"})
	if err != nil {
		t.Fatalf("Upsert conflict: %v", err)
	}
	if u2.ID != u.ID {
		t.Errorf("conflict created a new row: id %d, want %d", u2.ID, u.ID)
	}
	if u2.Login != renamed || u2.AccessToken != "gho_second" {
		t.Errorf("conflict did not update row: %+v", u2)
	}
}

func TestStore_GetByLoginNotFound(t *testing.T) {
	s := dbStore(t, nil)
	if _, err := s.GetByLogin(context.Background(), "definitely-not-a-real-login-xyz"); !errors.Is(err, ErrNotFound) {
		t.Errorf("GetByLogin error = %v, want ErrNotFound", err)
	}
}

// Two accounts can hold one login while a rename is pending; the current owner
// is the row updated most recently.
func TestStore_GetByLoginPrefersMostRecent(t *testing.T) {
	s := dbStore(t, nil)
	ctx := context.Background()
	login := "store-drift-login"

	stale := testdb.User(t, s.db, login)
	current := testdb.User(t, s.db, login)
	if _, err := s.db.Exec(ctx,
		`UPDATE users SET updated_at = now() - interval '30 days' WHERE id = $1`, stale); err != nil {
		t.Fatalf("age stale row: %v", err)
	}

	got, err := s.GetByLogin(ctx, login)
	if err != nil {
		t.Fatalf("GetByLogin: %v", err)
	}
	if got.ID != current {
		t.Errorf("GetByLogin returned id %d, want %d (most recently updated)", got.ID, current)
	}
}

func TestStore_TokenEncryptedAtRest(t *testing.T) {
	key := bytes.Repeat([]byte{7}, 32)
	s := dbStore(t, key)
	ctx := context.Background()

	u, err := s.Upsert(ctx, &User{GitHubID: 8100000002, Login: "store-crypto", AccessToken: "gho_secret"})
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	t.Cleanup(func() { _, _ = s.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, u.ID) })

	var stored string
	if err := s.db.QueryRow(ctx, `SELECT access_token FROM users WHERE id = $1`, u.ID).Scan(&stored); err != nil {
		t.Fatalf("read raw token: %v", err)
	}
	if stored == "gho_secret" {
		t.Error("token stored in plaintext")
	}

	got, err := s.GetByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.AccessToken != "gho_secret" {
		t.Errorf("AccessToken = %q, want the decrypted value", got.AccessToken)
	}
}

func TestStore_ListAll(t *testing.T) {
	s := dbStore(t, nil)
	id := testdb.User(t, s.db, "store-listall")

	all, err := s.ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	for _, u := range all {
		if u.ID == id {
			return
		}
	}
	t.Errorf("ListAll did not include user %d", id)
}
