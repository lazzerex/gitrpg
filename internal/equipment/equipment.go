package equipment

import (
	"sort"

	"github.com/lazzerex/gitrpg/internal/github"
)

type Item struct {
	Slug string
	Name string
	Icon string
	Lang string
}

var items = []Item{
	{Slug: "go-compiler", Name: "Go Compiler", Icon: "sword", Lang: "Go"},
	{Slug: "rust-axe", Name: "Rust Axe", Icon: "axe", Lang: "Rust"},
	{Slug: "typescript-lens", Name: "TypeScript Lens", Icon: "telescope", Lang: "TypeScript"},
	{Slug: "javascript-dagger", Name: "JavaScript Dagger", Icon: "sword", Lang: "JavaScript"},
	{Slug: "python-flask", Name: "Python Flask", Icon: "flask-conical", Lang: "Python"},
	{Slug: "csharp-bulwark", Name: "C# Bulwark", Icon: "shield", Lang: "C#"},
	{Slug: "java-hammer", Name: "Java Hammer", Icon: "hammer", Lang: "Java"},
	{Slug: "cpp-gauntlet", Name: "C++ Gauntlet", Icon: "shield-half", Lang: "C++"},
}

var byLang = func() map[string]Item {
	m := make(map[string]Item, len(items))
	for _, it := range items {
		m[it.Lang] = it
	}
	return m
}()

var bySlug = func() map[string]Item {
	m := make(map[string]Item, len(items))
	for _, it := range items {
		m[it.Slug] = it
	}
	return m
}()

// Lookup returns the catalog item for a slug.
func Lookup(slug string) (Item, bool) {
	it, ok := bySlug[slug]
	return it, ok
}

type Loadout struct {
	Weapon    *Item
	Shield    *Item
	Accessory *Item
}

type Slot struct {
	Label string
	Item  *Item
}

func (l Loadout) Slots() []Slot {
	return []Slot{
		{Label: "Weapon", Item: l.Weapon},
		{Label: "Shield", Item: l.Shield},
		{Label: "Accessory", Item: l.Accessory},
	}
}

func (l Loadout) Any() bool {
	return l.Weapon != nil || l.Shield != nil || l.Accessory != nil
}

// Top 3 distinct recognized languages by commit count fill Weapon/Shield/Accessory
// in rank order, not by tech-to-slot type (see GAME_DESIGN.md).
func Evaluate(s *github.Stats) Loadout {
	var loadout Loadout
	slots := [...]**Item{&loadout.Weapon, &loadout.Shield, &loadout.Accessory}

	i := 0
	for _, lang := range rankLanguages(s.Languages) {
		if i >= len(slots) {
			break
		}
		item, ok := byLang[lang]
		if !ok {
			continue
		}
		*slots[i] = &item
		i++
	}
	return loadout
}

// Alphabetical tie-break keeps this deterministic despite Go's randomized map order.
func rankLanguages(langs map[string]int) []string {
	names := make([]string, 0, len(langs))
	for lang := range langs {
		names = append(names, lang)
	}
	sort.Slice(names, func(i, j int) bool {
		if langs[names[i]] != langs[names[j]] {
			return langs[names[i]] > langs[names[j]]
		}
		return names[i] < names[j]
	})
	return names
}
