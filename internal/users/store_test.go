package users

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/lazzerex/gitrpg/internal/crypto"
)

func newTestStore(t *testing.T, key []byte) (*Store, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
	return NewStore(nil, key, logger), &buf
}

func TestDecryptToken_RoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	s, logs := newTestStore(t, key)

	sealed, err := crypto.Seal([]byte("gho_realtoken"), key)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if got := s.decryptToken(sealed); got != "gho_realtoken" {
		t.Errorf("decryptToken = %q, want %q", got, "gho_realtoken")
	}
	if logs.Len() != 0 {
		t.Errorf("unexpected log output: %s", logs)
	}
}

func TestDecryptToken_LegacyPlaintextIsSilent(t *testing.T) {
	s, logs := newTestStore(t, bytes.Repeat([]byte{1}, 32))

	if got := s.decryptToken("gho_legacyplaintext"); got != "gho_legacyplaintext" {
		t.Errorf("decryptToken = %q, want the value unchanged", got)
	}
	if logs.Len() != 0 {
		t.Errorf("legacy plaintext must not warn, got: %s", logs)
	}
}

func TestDecryptToken_KeyMismatchWarns(t *testing.T) {
	oldKey := bytes.Repeat([]byte{1}, 32)
	newKey := bytes.Repeat([]byte{2}, 32)

	sealed, err := crypto.Seal([]byte("gho_realtoken"), oldKey)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	s, logs := newTestStore(t, newKey)
	if got := s.decryptToken(sealed); got != sealed {
		t.Errorf("decryptToken = %q, want the stored value returned unchanged", got)
	}
	if !strings.Contains(logs.String(), "failed to decrypt") {
		t.Errorf("key mismatch must warn, got: %q", logs.String())
	}
}

func TestDecryptToken_NoKeyIsPassthrough(t *testing.T) {
	s, logs := newTestStore(t, nil)
	if got := s.decryptToken("gho_plain"); got != "gho_plain" {
		t.Errorf("decryptToken = %q, want passthrough", got)
	}
	if logs.Len() != 0 {
		t.Errorf("no-key mode must not warn, got: %s", logs)
	}
}
