package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionRoundtrip(t *testing.T) {
	secret := "test-secret"
	userID := int64(42)

	w := httptest.NewRecorder()
	if err := setSession(w, secret, userID, 3600, false); err != nil {
		t.Fatalf("setSession: %v", err)
	}

	r := &http.Request{Header: http.Header{"Cookie": {w.Header().Get("Set-Cookie")}}}
	got, err := getSession(r, secret)
	if err != nil {
		t.Fatalf("getSession: %v", err)
	}
	if got != userID {
		t.Fatalf("want %d, got %d", userID, got)
	}
}

func TestSessionWrongSecret(t *testing.T) {
	w := httptest.NewRecorder()
	setSession(w, "secret-a", 1, 3600, false)

	r := &http.Request{Header: http.Header{"Cookie": {w.Header().Get("Set-Cookie")}}}
	if _, err := getSession(r, "secret-b"); err == nil {
		t.Fatal("expected error with wrong secret")
	}
}
