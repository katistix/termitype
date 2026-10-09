// Package theme holds the monkeytype color themes, fuzzy search over them,
// and persistence of the selected theme.
//
// The theme table in themes_gen.go is generated from monkeytype's
// frontend/src/ts/constants/themes.ts.
package theme

import (
	"sort"
	"strings"
)

// Default is the theme used when none has been saved.
const Default = "serika_dark"

// Theme mirrors monkeytype's color variables.
type Theme struct {
	Name       string // id, e.g. "serika_dark"
	Bg         string // background
	Main       string // accent: title, stat values
	Caret      string // cursor
	Sub        string // untyped text, hints
	SubAlt     string // secondary background: borders
	Text       string // correctly typed text
	Error      string // wrong letters
	ErrorExtra string // extra letters typed past the word
}

// Display returns the name as monkeytype shows it: "serika dark".
func (t Theme) Display() string { return strings.ReplaceAll(t.Name, "_", " ") }

// All returns every theme, sorted by name.
func All() []Theme { return all }

// Get looks up a theme by name. Spaces and underscores are interchangeable
// and matching ignores case.
func Get(name string) (Theme, bool) {
	key := normalize(name)
	for _, t := range all {
		if t.Name == key {
			return t, true
		}
	}
	return Theme{}, false
}

// MustGet returns the named theme, falling back to Default.
func MustGet(name string) Theme {
	if t, ok := Get(name); ok {
		return t
	}
	t, _ := Get(Default)
	return t
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.ReplaceAll(s, " ", "_")
}

func init() {
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
}
