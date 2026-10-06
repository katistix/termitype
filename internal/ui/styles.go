// Package ui contains the bubbletea model and lipgloss styles.
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Palette: muted background, monkeytype-ish yellow accent.
	accent   = lipgloss.Color("#e2b714")
	muted    = lipgloss.Color("#646669")
	fg       = lipgloss.Color("#d1d0c5")
	bad      = lipgloss.Color("#ca4754")
	bg       = lipgloss.Color("#323437")
	cursorBg = lipgloss.Color("#e2b714")
	cursorFg = lipgloss.Color("#323437")
)

type Styles struct {
	App     lipgloss.Style
	Title   lipgloss.Style
	Hint    lipgloss.Style
	Box     lipgloss.Style
	Pending lipgloss.Style
	Correct lipgloss.Style
	Wrong   lipgloss.Style
	Extra   lipgloss.Style
	Cursor  lipgloss.Style
	Stat    lipgloss.Style
	StatVal lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		App:     lipgloss.NewStyle().Padding(1, 2),
		Title:   lipgloss.NewStyle().Bold(true).Foreground(accent),
		Hint:    lipgloss.NewStyle().Foreground(muted),
		Box:     lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(bg),
		Pending: lipgloss.NewStyle().Foreground(muted),
		Correct: lipgloss.NewStyle().Foreground(fg),
		Wrong:   lipgloss.NewStyle().Foreground(bad),
		Extra:   lipgloss.NewStyle().Foreground(bad).Background(lipgloss.Color("#4a2b2e")),
		Cursor:  lipgloss.NewStyle().Background(cursorBg).Foreground(cursorFg),
		Stat:    lipgloss.NewStyle().Foreground(muted),
		StatVal: lipgloss.NewStyle().Bold(true).Foreground(accent),
	}
}
