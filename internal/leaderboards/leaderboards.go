package leaderboards

type Entry struct {
	Rank      int
	Login     string
	AvatarURL string
	Class     string
	Level     int
	TotalXP   int
}

const PageSize = 50

type Page struct {
	Entries    []Entry
	Page       int
	TotalPages int
}

func (p Page) HasPrev() bool { return p.Page > 1 }
func (p Page) HasNext() bool { return p.Page < p.TotalPages }
