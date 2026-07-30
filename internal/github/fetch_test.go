package github

import (
	"testing"
	"time"
)

func TestYearWindows(t *testing.T) {
	createdAt := time.Date(2019, 6, 15, 10, 30, 0, 0, time.UTC)
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)

	ws := yearWindows(createdAt, now)
	if len(ws) == 0 {
		t.Fatal("no windows returned")
	}

	first := ws[0]
	if want := time.Date(2019, 6, 15, 0, 0, 0, 0, time.UTC); !first.from.Equal(want) {
		t.Errorf("first window starts %v, want %v (day-aligned createdAt)", first.from, want)
	}
	if last := ws[len(ws)-1]; !last.to.Equal(now) {
		t.Errorf("last window ends %v, want %v", last.to, now)
	}

	for i, w := range ws {
		if !w.from.Before(w.to) {
			t.Errorf("window %d: from %v not before to %v", i, w.from, w.to)
		}
		if limit := w.from.AddDate(1, 0, 0); w.to.After(limit) {
			t.Errorf("window %d spans more than one year: %v → %v", i, w.from, w.to)
		}
		if i > 0 {
			if want := ws[i-1].to.Add(time.Second); !w.from.Equal(want) {
				t.Errorf("window %d: from %v, want %v (disjoint from previous)", i, w.from, want)
			}
		}
	}
}

func TestYearWindows_LeapDayCreatedAt(t *testing.T) {
	createdAt := time.Date(2020, 2, 29, 8, 0, 0, 0, time.UTC)
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)

	for i, w := range yearWindows(createdAt, now) {
		if limit := w.from.AddDate(1, 0, 0); w.to.After(limit) {
			t.Errorf("window %d spans more than one year: %v → %v", i, w.from, w.to)
		}
	}
}

func TestYearWindows_AccountYoungerThanOneYear(t *testing.T) {
	createdAt := time.Date(2026, 3, 1, 9, 0, 0, 0, time.UTC)
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)

	ws := yearWindows(createdAt, now)
	if len(ws) != 1 {
		t.Fatalf("got %d windows, want 1", len(ws))
	}
	if !ws[0].to.Equal(now) {
		t.Errorf("window ends %v, want %v", ws[0].to, now)
	}
}
