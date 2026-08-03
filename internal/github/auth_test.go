package github

import (
	"errors"
	"net/http"
	"testing"
)

func TestStatusError_UnauthorizedIsDetectable(t *testing.T) {
	err := statusError(http.StatusUnauthorized)
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("statusError(401) = %v, want it to wrap ErrUnauthorized", err)
	}
	if errors.Is(statusError(http.StatusBadGateway), ErrUnauthorized) {
		t.Error("statusError(502) must not wrap ErrUnauthorized")
	}
}

func TestSyncStatusFor(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"success", nil, "success"},
		{"unauthorized", statusError(http.StatusUnauthorized), StatusUnauthorized},
		{"wrapped unauthorized", errors.Join(errors.New("sync aborted"), ErrUnauthorized), StatusUnauthorized},
		{"other failure", errors.New("boom"), "failed"},
		{"other http status", statusError(http.StatusBadGateway), "failed"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := syncStatusFor(c.err); got != c.want {
				t.Errorf("syncStatusFor(%v) = %q, want %q", c.err, got, c.want)
			}
		})
	}
}
