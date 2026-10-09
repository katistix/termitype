// Package ui contains the bubbletea models and lipgloss styles.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/katistix/termitype/internal/theme"
)

// Styles are built from a theme. Every style sets the theme background:
// ANSI resets between styled segments would otherwise punch holes in it.
type Styles struct {
	Bg      lipgloss.Color
	Base    lipgloss.Style // plain text and gaps on the theme background
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

func NewStyles(t theme.Theme) Styles {
	bg := lipgloss.Color(t.Bg)
	base := lipgloss.NewStyle().Background(bg)
	fg := func(c string) lipgloss.Style { return base.Foreground(lipgloss.Color(c)) }
	return Styles{
		Bg:      bg,
		Base:    base,
		App:     base.Padding(1, 2),
		Title:   fg(t.Main).Bold(true),
		Hint:    fg(t.Sub),
		Box:     base.Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(t.SubAlt)).BorderBackground(bg),
		Pending: fg(t.Sub),
		Correct: fg(t.Text),
		Wrong:   fg(t.Error),
		Extra:   fg(t.ErrorExtra),
		Cursor:  lipgloss.NewStyle().Background(lipgloss.Color(t.Caret)).Foreground(bg),
		Stat:    fg(t.Sub),
		StatVal: fg(t.Main).Bold(true),
	}
}

// fill paints content over the whole terminal in the theme background.
func (s Styles) fill(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return content
	}
	// Joins and padding emit plain spaces after a style reset; re-apply
	// the background after each reset and at each line start.
	bgSeq, _, _ := strings.Cut(s.Base.Render("x"), "x")
	if bgSeq != "" {
		content = strings.ReplaceAll(content, "\x1b[0m", "\x1b[0m"+bgSeq)
		content = bgSeq + strings.ReplaceAll(content, "\n", "\n"+bgSeq)
	}
	return lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, content,
		lipgloss.WithWhitespaceBackground(s.Bg))
}
