package quests

import (
	"testing"
	"time"
)

func TestWeekStart(t *testing.T) {
	tests := []struct {
		in   time.Time
		want time.Time
	}{
		{time.Date(2026, 7, 31, 15, 30, 0, 0, time.UTC), time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC)}, // Friday
		{time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC), time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC)},   // Monday itself
		{time.Date(2026, 8, 2, 23, 59, 0, 0, time.UTC), time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC)},  // Sunday
		{time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), time.Date(2025, 12, 29, 0, 0, 0, 0, time.UTC)},  // year boundary
	}
	for _, tt := range tests {
		if got := weekStart(tt.in); !got.Equal(tt.want) {
			t.Errorf("weekStart(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestProgressOf(t *testing.T) {
	tests := []struct {
		current, baseline, target, want int
	}{
		{110, 100, 20, 10},
		{130, 100, 20, 20}, // overshoot clamps to target
		{95, 100, 20, 0},   // corrected counts clamp to zero
		{100, 100, 20, 0},
	}
	for _, tt := range tests {
		if got := progressOf(tt.current, tt.baseline, tt.target); got != tt.want {
			t.Errorf("progressOf(%d, %d, %d) = %d, want %d", tt.current, tt.baseline, tt.target, got, tt.want)
		}
	}
}
