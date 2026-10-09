package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/katistix/termitype/internal/config"
	"github.com/katistix/termitype/internal/theme"
	"github.com/katistix/termitype/internal/typing"
	"github.com/katistix/termitype/internal/words"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Every(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model is the bubbletea model for the typing test.
type Model struct {
	cfg    config.Config
	engine *typing.Engine
	styles Styles
	width  int
	height int
}

// New creates a model with a fresh test from cfg, drawn in theme t.
func New(cfg config.Config, t theme.Theme) Model {
	return Model{
		cfg:    cfg,
		engine: newEngine(cfg),
		styles: NewStyles(t),
	}
}

// newEngine builds a typing engine matching the mode.
// Time mode pre-generates a generous buffer and refills as you type.
func newEngine(cfg config.Config) *typing.Engine {
	if cfg.Mode == config.ModeWords {
		tcfg := typing.Config{WordCount: cfg.Words}
		return typing.New(words.Generate(cfg.Words), tcfg)
	}
	tcfg := typing.Config{TimeLimit: time.Duration(cfg.Seconds) * time.Second}
	return typing.New(words.Generate(timeBufferWords(cfg.Seconds)), tcfg)
}

// timeBufferWords estimates enough words so refills are rare.
func timeBufferWords(seconds int) int {
	n := seconds * 4 // ~48 wpm in 5-char words, with headroom
	if n < 60 {
		n = 60
	}
	return n
}

func (m *Model) reset() {
	m.engine = newEngine(m.cfg)
}

func (m Model) Init() tea.Cmd { return tick() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		if m.engine.CheckTimeout() {
			return m, nil // stopped: no more ticks needed
		}
		if m.engine.Done() {
			return m, nil
		}
		return m, tick()

	case tea.KeyMsg:
		// Option+Backspace arrives as Alt+Backspace on macOS terminals.
		// Treat it (and ctrl+w) as delete-word, like monkeytype.
		if msg.Type == tea.KeyBackspace && msg.Alt {
			m.engine.DeleteWord()
			return m, nil
		}
		// Some terminals send ctrl+backspace as ctrl+w.
		if s := msg.String(); s == "ctrl+backspace" || s == "alt+backspace" {
			m.engine.DeleteWord()
			return m, nil
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			// Esc restarts with fresh words.
			m.reset()
			return m, tick()

		case tea.KeyTab, tea.KeyEnter:
			return m, nil

		case tea.KeyBackspace:
			m.engine.Backspace()
			return m, nil

		case tea.KeyCtrlW:
			m.engine.DeleteWord()
			return m, nil

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.typeRune(r)
			}
			return m, nil

		case tea.KeySpace:
			m.typeRune(' ')
			return m, nil
		}
	}
	return m, nil
}

// typeRune types one rune, refilling words in time mode and enforcing timeout.
func (m *Model) typeRune(r rune) {
	if m.engine.Done() {
		return
	}
	m.engine.TypeRune(r)
	if m.cfg.Mode == config.ModeTime && m.engine.NeedsMoreWords(30) {
		m.engine.AppendWords(words.GenerateSlice(30))
	}
	m.engine.CheckTimeout()
}

func (m Model) View() string {
	s := m.styles

	var b strings.Builder
	b.WriteString(s.Title.Render("termitype"))
	b.WriteString(s.Hint.Render("  " + m.cfg.Label()))
	b.WriteString("\n\n")

	if m.engine.Done() {
		b.WriteString(m.renderResults())
	} else {
		b.WriteString(m.renderLiveStats())
		b.WriteString("\n\n")
		b.WriteString(m.renderText())
	}

	return s.fill(s.App.Render(b.String()), m.width, m.height)
}

func (m Model) renderLiveStats() string {
	s := m.styles
	parts := ""
	if m.cfg.Mode == config.ModeTime {
		secs := int(m.engine.Remaining().Round(time.Second).Seconds())
		parts += s.Stat.Render("time ") + s.StatVal.Render(fmt.Sprintf("%d", secs)) + s.Stat.Render("   ")
	}
	wpm := fmt.Sprintf("%.0f", m.engine.WPM())
	acc := fmt.Sprintf("%.0f%%", m.engine.Accuracy()*100)
	parts += s.Stat.Render("wpm ") + s.StatVal.Render(wpm) +
		s.Stat.Render("   acc ") + s.StatVal.Render(acc)
	return parts
}

// renderText draws words with per-letter states. The caret stays on the
// current word; wrong letters show as extras until Space submits the word.
func (m Model) renderText() string {
	s := m.styles
	e := m.engine

	// Window the view so long time-mode buffers stay readable.
	start, end := 0, e.WordCount()
	if e.WordCount() > 60 {
		start = e.CurrentWord() - 10
		if start < 0 {
			start = 0
		}
		end = start + 60
		if end > e.WordCount() {
			end = e.WordCount()
			start = end - 60
		}
	}

	var sb strings.Builder
	for i := start; i < end; i++ {
		tw := []rune(e.TargetWord(i))
		in := []rune(e.InputWord(i))
		caret := i == e.CurrentWord() && !e.Done()
		caretPos := len(in) // caret sits after last typed rune

		// Target letters + extras.
		max := len(tw)
		if len(in) > max {
			max = len(in)
		}
		for j := 0; j < max; j++ {
			var ch string
			var st int
			isExtra := j >= len(tw)
			if isExtra {
				ch = string(in[j])
			} else {
				ch = string(tw[j])
			}
			if j < len(in) {
				st = e.LetterStatus(i, j)
			} else {
				st = 0
			}
			atCaret := caret && j == caretPos
			switch {
			case atCaret && isExtra:
				sb.WriteString(s.Cursor.Render(ch))
			case atCaret:
				sb.WriteString(s.Cursor.Render(ch))
			case j < len(in) && st == 1:
				sb.WriteString(s.Correct.Render(ch))
			case j < len(in) && st == 2:
				if isExtra {
					sb.WriteString(s.Extra.Render(ch))
				} else {
					sb.WriteString(s.Wrong.Render(ch))
				}
			default:
				sb.WriteString(s.Pending.Render(ch))
			}
		}
		// Caret on the inter-word gap (input longer than nothing to show,
		// or empty caret past last letter handled above for last word).
		if caret && caretPos == len(tw) && len(in) == len(tw) {
			// Caret is on the space after this word (or end marker).
			if i+1 < end {
				sb.WriteString(s.Cursor.Render(" "))
			} else {
				sb.WriteString(s.Cursor.Render("·"))
			}
		} else if i+1 < end {
			sb.WriteString(s.Base.Render(" "))
		}
	}
	w := m.width - 8
	if w <= 0 || w > 80 {
		w = 80
	}
	return s.Box.Width(w).Render(sb.String())
}

func (m Model) renderResults() string {
	s := m.styles
	wpm := fmt.Sprintf("%.0f", m.engine.WPM())
	raw := fmt.Sprintf("%.0f", m.engine.RawWPM())
	acc := fmt.Sprintf("%.1f%%", m.engine.Accuracy()*100)
	correct := m.engine.CorrectChars()
	total := m.engine.TypedChars()
	lines := []string{
		s.Stat.Render(fmt.Sprintf("%s   ", m.cfg.Label())) + s.StatVal.Render("done"),
		s.Stat.Render("wpm  ") + s.StatVal.Render(wpm) + s.Stat.Render("   raw  ") + s.StatVal.Render(raw),
		s.Stat.Render("acc  ") + s.StatVal.Render(acc),
		s.Stat.Render("chars ") + s.StatVal.Render(fmt.Sprintf("%d/%d", correct, total)),
	}
	return s.Box.Render(strings.Join(lines, "\n"))
}
