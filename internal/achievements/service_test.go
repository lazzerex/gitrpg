package achievements

import (
	"reflect"
	"testing"
)

func TestDiffSlugs(t *testing.T) {
	tests := []struct {
		name     string
		earned   []string
		existing []string
		want     []string
	}{
		{"all new", []string{"a", "b"}, nil, []string{"a", "b"}},
		{"some new", []string{"a", "b", "c"}, []string{"a"}, []string{"b", "c"}},
		{"none new", []string{"a", "b"}, []string{"a", "b"}, nil},
		{"empty earned", nil, []string{"a"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diffSlugs(tt.earned, tt.existing); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diffSlugs(%v, %v) = %v, want %v", tt.earned, tt.existing, got, tt.want)
			}
		})
	}
}
