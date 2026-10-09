package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/katistix/termitype/internal/theme"
)

// Picker is a fuzzy theme selector. The whole screen previews the
// highlighted theme; each list row is drawn in its own theme's colors.
type Picker struct {
	query    string
	matches  []theme.Match
	cursor   int // index into matches
	offset   int // first visible match
	current  string
	width    int
	height   int
	chosen   *theme.Theme
	canceled bool
}

// NewPicker starts with the cursor on the current theme.
func NewPicker(current theme.Theme) Picker {
	p := Picker{current: current.Name}
	p.refilter()
	for i, mt := range p.matches {
		if mt.Theme.Name == current.Name {
			p.cursor = i
		}
	}
	return p
}

// Chosen returns the selected theme, or false if the picker was canceled.
func (p Picker) Chosen() (theme.Theme, bool) {
	if p.chosen == nil {
		return theme.Theme{}, false
	}
	return *p.chosen, true
}

func (p *Picker) refilter() {
	p.matches = theme.Filter(theme.All(), p.query)
	p.cursor, p.offset = 0, 0
}

func (p Picker) highlighted() theme.Theme {
	if len(p.matches) == 0 {
		return theme.MustGet(p.current)
	}
	return p.matches[p.cursor].Theme
}

func (p Picker) Init() tea.Cmd { return nil }

func (p Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width, p.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			p.canceled = true
			return p, tea.Quit
		case "enter":
			if len(p.matches) > 0 {
				t := p.matches[p.cursor].Theme
				p.chosen = &t
				return p, tea.Quit
			}
		case "up", "ctrl+p", "ctrl+k":
			p.move(-1)
		case "down", "ctrl+n", "ctrl+j", "tab":
			p.move(1)
		case "shift+tab":
			p.move(-1)
		case "pgup":
			p.move(-p.listHeight())
		case "pgdown":
			p.move(p.listHeight())
		case "home":
			p.move(-len(p.matches))
		case "end":
			p.move(len(p.matches))
		case "backspace":
			if r := []rune(p.query); len(r) > 0 {
				p.query = string(r[:len(r)-1])
				p.refilter()
			}
		case "ctrl+u", "ctrl+w", "alt+backspace":
			p.query = ""
			p.refilter()
		default:
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
				p.query += string(msg.Runes)
				p.refilter()
			}
		}
	}
	return p, nil
}

func (p *Picker) move(d int) {
	if len(p.matches) == 0 {
		return
	}
	p.cursor = max(0, min(len(p.matches)-1, p.cursor+d))
	h := p.listHeight()
	if p.cursor < p.offset {
		p.offset = p.cursor
	} else if p.cursor >= p.offset+h {
		p.offset = p.cursor - h + 1
	}
}

// listHeight is the number of theme rows that fit on screen.
func (p Picker) listHeight() int {
	// App padding (2) + header (2) + prompt (2) + footer (2).
	h := p.height - 8
	if p.height == 0 {
		h = 15
	}
	return max(1, h)
}

const rowWidth = 28

func (p Picker) View() string {
	t := p.highlighted()
	s := NewStyles(t)

	header := s.Title.Render("termitype") + s.Hint.Render("  theme")
	prompt := s.StatVal.Render("> ") + s.Correct.Render(p.query) + s.Cursor.Render(" ")

	list := p.renderList(s)
	body := list
	// Preview only when there's room beside the list.
	if p.width == 0 || p.width >= rowWidth+50 {
		gap := s.Base.Render("    ")
		body = lipgloss.JoinHorizontal(lipgloss.Top, list, gap, renderPreview(s, t))
	}

	footer := s.Hint.Render(fmt.Sprintf("%d/%d", len(p.matches), len(theme.All()))) +
		s.Hint.Render("   ↑/↓ move · enter select · esc cancel")

	view := lipgloss.JoinVertical(lipgloss.Left, header, "", prompt, "", body, "", footer)
	return s.fill(s.App.Render(view), p.width, p.height)
}

func (p Picker) renderList(s Styles) string {
	if len(p.matches) == 0 {
		return s.Hint.Width(rowWidth + 2).Render("no matching themes")
	}
	h := p.listHeight()
	end := min(len(p.matches), p.offset+h)
	rows := make([]string, 0, h)
	for i := p.offset; i < end; i++ {
		marker := s.Base.Render("  ")
		if i == p.cursor {
			marker = s.StatVal.Render("▸ ")
		}
		rows = append(rows, marker+renderRow(p.matches[i], p.matches[i].Theme.Name == p.current))
	}
	// Pad so the preview doesn't jump around as the list shrinks.
	for len(rows) < h {
		rows = append(rows, s.Base.Render(strings.Repeat(" ", rowWidth+2)))
	}
	return strings.Join(rows, "\n")
}

// renderRow draws a theme name on its own background with matched
// characters in its accent color, plus main/sub/text color dots.
func renderRow(mt theme.Match, isCurrent bool) string {
	t := mt.Theme
	base := lipgloss.NewStyle().Background(lipgloss.Color(t.Bg))
	name := base.Foreground(lipgloss.Color(t.Text))
	hit := base.Foreground(lipgloss.Color(t.Main)).Bold(true).Underline(true)

	matched := make(map[int]bool, len(mt.Indices))
	for _, i := range mt.Indices {
		matched[i] = true
	}

	label := t.Display()
	if isCurrent {
		label += " *"
	}
	// Truncate to leave room for padding and the color dots.
	maxName := rowWidth - 8
	if len(label) > maxName {
		label = label[:maxName-1] + "…"
	}

	var b strings.Builder
	b.WriteString(base.Render(" "))
	for i, r := range label {
		if matched[i] {
			b.WriteString(hit.Render(string(r)))
		} else {
			b.WriteString(name.Render(string(r)))
		}
	}
	pad := rowWidth - 1 - lipgloss.Width(label) - 6
	b.WriteString(base.Render(strings.Repeat(" ", max(0, pad))))
	for _, c := range []string{t.Main, t.Sub, t.Text} {
		b.WriteString(base.Foreground(lipgloss.Color(c)).Render("●"))
		b.WriteString(base.Render(" "))
	}
	return b.String()
}

// renderPreview draws a mock test mid-way through, showing every
// letter state: correct, wrong, extra, caret and pending.
func renderPreview(s Styles, t theme.Theme) string {
	type seg struct {
		st   lipgloss.Style
		text string
	}
	segs := []seg{
		{s.Correct, "the quick"}, {s.Base, " "},
		{s.Correct, "br"}, {s.Wrong, "w"}, {s.Correct, "wn"}, {s.Base, " "},
		{s.Correct, "fox"}, {s.Extra, "xx"}, {s.Base, " "},
		{s.Correct, "jum"}, {s.Cursor, "p"}, {s.Pending, "s"}, {s.Base, " "},
		{s.Pending, "over the lazy dog and then some more words to type"},
	}
	var text strings.Builder
	for _, sg := range segs {
		text.WriteString(sg.st.Render(sg.text))
	}
	stats := s.Stat.Render("time ") + s.StatVal.Render("12") +
		s.Stat.Render("   wpm ") + s.StatVal.Render("87") +
		s.Stat.Render("   acc ") + s.StatVal.Render("96%")
	title := s.Title.Render(t.Display())
	return lipgloss.JoinVertical(lipgloss.Left, title, "", stats, "", s.Box.Width(42).Render(text.String()))
}
