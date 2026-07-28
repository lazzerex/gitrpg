package leaderboards

import "testing"

func TestPaginate(t *testing.T) {
	cases := []struct {
		name           string
		total, page    int
		wantPage       int
		wantTotalPages int
	}{
		{"empty leaderboard", 0, 1, 1, 1},
		{"single page", 10, 1, 1, 1},
		{"exact multiple of page size", PageSize * 2, 2, 2, 2},
		{"partial last page", PageSize*2 + 1, 3, 3, 3},
		{"page below 1 clamps to 1", 100, 0, 1, 2},
		{"negative page clamps to 1", 100, -5, 1, 2},
		{"page beyond total clamps to last", PageSize, 99, 1, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotPage, gotTotalPages := paginate(c.total, c.page)
			if gotPage != c.wantPage || gotTotalPages != c.wantTotalPages {
				t.Fatalf("paginate(%d, %d) = (%d, %d), want (%d, %d)",
					c.total, c.page, gotPage, gotTotalPages, c.wantPage, c.wantTotalPages)
			}
		})
	}
}

func TestPage_HasPrevHasNext(t *testing.T) {
	p := Page{Page: 2, TotalPages: 3}
	if !p.HasPrev() {
		t.Fatal("expected HasPrev true")
	}
	if !p.HasNext() {
		t.Fatal("expected HasNext true")
	}

	last := Page{Page: 3, TotalPages: 3}
	if last.HasNext() {
		t.Fatal("expected HasNext false on last page")
	}

	first := Page{Page: 1, TotalPages: 3}
	if first.HasPrev() {
		t.Fatal("expected HasPrev false on first page")
	}
}
