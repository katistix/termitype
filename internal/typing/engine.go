// Package typing holds the core typing-test state machine.
// It mirrors monkeytype: typing stays on the current word (extra letters
// are kept as extras), Space submits the word and advances, Backspace
// edits within the word, and you can only jump back to fix incorrect words.
package typing

import (
	"strings"
	"time"
)

// maxExtra caps extra letters beyond the target word length,
// like monkeytype which stops accepting endless extras.
const maxExtra = 10

// Config controls a test.
type Config struct {
	WordCount int           // words per test
	TimeLimit time.Duration // 0 = no limit, finish when all words done
}

// DefaultConfig is words-mode with 30 words (time mode sets TimeLimit).
func DefaultConfig() Config {
	return Config{WordCount: 30, TimeLimit: 0}
}

// Engine tracks target words vs per-word user input.
type Engine struct {
	cfg     Config
	target  []string
	input   [][]rune // typed runes per word, len == len(target)
	current int      // index of active word
	started time.Time
	ended   *time.Time
	begun   bool

	// Keystroke counters for monkeytype-style accuracy.
	strokes, correctStrokes int
}

// New creates an Engine for the given target text (space-separated words).
func New(target string, cfg Config) *Engine {
	ws := strings.Fields(target)
	in := make([][]rune, len(ws))
	return &Engine{cfg: cfg, target: ws, input: in}
}

// Target returns the full target text.
func (e *Engine) Target() string { return strings.Join(e.target, " ") }

// Typed returns submitted + current input joined as text (no trailing space).
func (e *Engine) Typed() string {
	end := e.current
	if e.ended != nil && end < len(e.input) {
		end++
	}
	parts := make([]string, 0, end)
	for i := 0; i < end && i < len(e.input); i++ {
		parts = append(parts, string(e.input[i]))
	}
	return strings.Join(parts, " ")
}

func (e *Engine) wordSubmitted(i int) bool { return i < e.current }

// WordCount returns number of target words.
func (e *Engine) WordCount() int { return len(e.target) }

// TargetWord returns target word i.
func (e *Engine) TargetWord(i int) string {
	if i < 0 || i >= len(e.target) {
		return ""
	}
	return e.target[i]
}

// InputWord returns typed input for word i.
func (e *Engine) InputWord(i int) string {
	if i < 0 || i >= len(e.input) {
		return ""
	}
	return string(e.input[i])
}

// CurrentWord is the active word index.
func (e *Engine) CurrentWord() int { return e.current }

// CursorPos is the caret offset inside the current word.
func (e *Engine) CursorPos() int {
	if e.current >= len(e.input) {
		return 0
	}
	return len(e.input[e.current])
}

// Done reports whether the test is complete.
func (e *Engine) Done() bool { return e.ended != nil }

// Started reports whether the user typed anything yet.
func (e *Engine) Started() bool { return e.begun }

func (e *Engine) ensureStarted() {
	if !e.begun {
		e.begun = true
		e.started = time.Now()
	}
}

// TypeRune handles a letter (non-space) keystroke.
// It appends to the current word and never auto-advances:
// wrong letters become extras you must fix or leave until Space.
func (e *Engine) TypeRune(r rune) bool {
	if e.Done() {
		return false
	}
	if r == '\n' || r == '\r' || r == '\t' {
		return true
	}
	if r == ' ' {
		e.Space()
		return true
	}
	e.ensureStarted()
	if e.current >= len(e.target) {
		return true
	}
	tw := []rune(e.target[e.current])
	cur := e.input[e.current]
	if len(cur) >= len(tw)+maxExtra {
		return true // cap extras like monkeytype
	}
	e.input[e.current] = append(cur, r)
	e.strokes++
	if len(cur) < len(tw) && r == tw[len(cur)] {
		e.correctStrokes++
	}
	e.checkWordsDone()
	return true
}

// Space submits the current word and advances, like monkeytype.
// Empty input ignores the space (no double spaces). A correct word's
// trailing space counts as a correct keystroke.
func (e *Engine) Space() {
	if e.Done() {
		return
	}
	if e.current >= len(e.target) {
		return
	}
	if len(e.input[e.current]) == 0 {
		return // ignore leading / double space
	}
	e.ensureStarted()
	e.strokes++
	if string(e.input[e.current]) == e.target[e.current] {
		e.correctStrokes++
	}
	e.current++
	if e.current >= len(e.target) {
		e.Finish()
	}
}

func (e *Engine) checkWordsDone() {
	// Words mode finishes when the last word exactly matches,
	// no trailing space required.
	if e.cfg.TimeLimit == 0 && len(e.target) > 0 {
		last := len(e.target) - 1
		if e.current == last && string(e.input[last]) == e.target[last] {
			e.Finish()
		}
	}
}

// Backspace deletes one char in the current word. At word start it jumps
// back to the previous word only if that word is incorrect (monkeytype
// won't let you edit already-correct words).
func (e *Engine) Backspace() {
	if e.Done() || e.current >= len(e.input) {
		return
	}
	if len(e.input[e.current]) > 0 {
		e.input[e.current] = e.input[e.current][:len(e.input[e.current])-1]
		return
	}
	if e.current == 0 {
		return
	}
	prev := e.current - 1
	if string(e.input[prev]) == e.target[prev] {
		return // correct words are locked
	}
	e.current = prev
}

// DeleteWord clears the current word (ctrl+w / option+backspace).
// If the current word is already empty it jumps back and clears the
// previous incorrect word.
func (e *Engine) DeleteWord() {
	if e.Done() {
		return
	}
	if e.current >= len(e.input) {
		return
	}
	if len(e.input[e.current]) > 0 {
		e.input[e.current] = nil
		return
	}
	if e.current == 0 {
		return
	}
	prev := e.current - 1
	if string(e.input[prev]) == e.target[prev] {
		return // don't wipe correct words
	}
	e.current = prev
	e.input[prev] = nil
}

// BackspaceWord is an alias kept for callers (ctrl+w / alt+backspace).
func (e *Engine) BackspaceWord() { e.DeleteWord() }

// Reset clears input and timers, keeping the same words.
func (e *Engine) Reset() {
	for i := range e.input {
		e.input[i] = nil
	}
	e.current = 0
	e.begun = false
	e.ended = nil
	e.strokes, e.correctStrokes = 0, 0
}

// SetTarget replaces the words and resets.
func (e *Engine) SetTarget(target string) {
	ws := strings.Fields(target)
	e.target = ws
	e.input = make([][]rune, len(ws))
	e.current = 0
	e.Reset()
}

// AppendWords adds more target words (time mode refill).
func (e *Engine) AppendWords(ws []string) {
	if len(ws) == 0 {
		return
	}
	e.target = append(e.target, ws...)
	for range ws {
		e.input = append(e.input, nil)
	}
}

// AppendText extends the target (splits on whitespace).
func (e *Engine) AppendText(suffix string) { e.AppendWords(strings.Fields(suffix)) }

// RemainingWords counts words from current to end.
func (e *Engine) RemainingWords() int {
	if e.current >= len(e.target) {
		return 0
	}
	return len(e.target) - e.current
}

// NeedsMoreWords reports whether fewer than threshold words remain.
func (e *Engine) NeedsMoreWords(threshold int) bool {
	return e.RemainingWords() < threshold
}

// Finish marks the test complete.
func (e *Engine) Finish() {
	if e.ended == nil {
		t := time.Now()
		e.ended = &t
	}
}

// Elapsed is capped at TimeLimit so tick overshoot doesn't deflate WPM.
func (e *Engine) Elapsed() time.Duration {
	if !e.begun {
		return 0
	}
	var d time.Duration
	if e.ended != nil {
		d = e.ended.Sub(e.started)
	} else {
		d = time.Since(e.started)
	}
	if e.cfg.TimeLimit > 0 && d > e.cfg.TimeLimit {
		return e.cfg.TimeLimit
	}
	return d
}

// TimeLimit returns the configured limit (0 = none).
func (e *Engine) TimeLimit() time.Duration { return e.cfg.TimeLimit }

// Remaining returns time left in time mode.
func (e *Engine) Remaining() time.Duration {
	if e.cfg.TimeLimit <= 0 || !e.begun {
		if e.cfg.TimeLimit > 0 {
			return e.cfg.TimeLimit
		}
		return 0
	}
	rem := e.cfg.TimeLimit - e.Elapsed()
	if rem < 0 {
		return 0
	}
	return rem
}

// CheckTimeout finishes the test if the time limit expired.
func (e *Engine) CheckTimeout() bool {
	if e.ended != nil {
		return true
	}
	if e.cfg.TimeLimit > 0 && e.begun && time.Since(e.started) >= e.cfg.TimeLimit {
		e.Finish()
		return true
	}
	return false
}

// CorrectChars counts correct letters in typed words plus one per
// correctly submitted word (its trailing space), monkeytype-style.
func (e *Engine) CorrectChars() int {
	n := 0
	end := e.current
	if e.ended != nil && end < len(e.input) {
		end = e.current + 1
	}
	for i := 0; i < end && i < len(e.target); i++ {
		tw := []rune(e.target[i])
		in := e.input[i]
		for j, r := range in {
			if j < len(tw) && r == tw[j] {
				n++
			}
		}
		// Submitted (or finished) exact words earn their space.
		if i < e.current && string(in) == e.target[i] {
			n++
		}
		if e.ended != nil && i == e.current && string(in) == e.target[i] && e.cfg.TimeLimit == 0 {
			// Last word finished without trailing space: no extra point.
		}
	}
	return n
}

// TypedChars counts all typed letters plus submitted spaces (for raw WPM).
func (e *Engine) TypedChars() int {
	n := 0
	end := e.current
	if e.ended != nil && end < len(e.input) {
		end = e.current + 1
	}
	for i := 0; i < end && i < len(e.input); i++ {
		n += len(e.input[i])
		if i < e.current {
			n++ // space that submitted the word
		}
	}
	return n
}

// Accuracy is correct keystrokes / total keystrokes (monkeytype-style).
// Corrections count: fixing a mistake still leaves the wrong keystroke.
func (e *Engine) Accuracy() float64 {
	if e.strokes == 0 {
		return 1
	}
	return float64(e.correctStrokes) / float64(e.strokes)
}

// WPM = (correct chars / 5) / minutes.
func (e *Engine) WPM() float64 {
	mins := e.Elapsed().Minutes()
	if mins <= 0 {
		return 0
	}
	return (float64(e.CorrectChars()) / 5) / mins
}

// RawWPM counts all typed chars regardless of correctness.
func (e *Engine) RawWPM() float64 {
	mins := e.Elapsed().Minutes()
	if mins <= 0 {
		return 0
	}
	return (float64(e.TypedChars()) / 5) / mins
}

// LetterStatus reports per-letter state of word i, position j:
// 0 = pending, 1 = correct, 2 = incorrect/extra.
func (e *Engine) LetterStatus(i, j int) int {
	if i < 0 || i >= len(e.target) {
		return 0
	}
	tw := []rune(e.target[i])
	var in []rune
	if i < len(e.input) {
		in = e.input[i]
	}
	if j >= len(in) {
		return 0
	}
	if j < len(tw) && in[j] == tw[j] {
		return 1
	}
	return 2
}

// WordError reports whether word i was submitted (or finished) incorrectly.
func (e *Engine) WordError(i int) bool {
	if i >= len(e.target) {
		return false
	}
	if i < e.current || (e.ended != nil && i <= e.current) {
		return string(e.input[i]) != e.target[i]
	}
	return false
}
