package typing

import "testing"

func typeStr(e *Engine, s string) {
	for _, r := range s {
		e.TypeRune(r)
	}
}

func TestStaysOnWordUntilSpace(t *testing.T) {
	e := New("hello world", DefaultConfig())
	typeStr(e, "hellx")
	if e.CurrentWord() != 0 {
		t.Fatalf("word=%d want 0 (must not advance on wrong letter)", e.CurrentWord())
	}
	if got := e.InputWord(0); got != "hellx" {
		t.Fatalf("input=%q want hellx (extra kept)", got)
	}
	e.Space()
	if e.CurrentWord() != 1 {
		t.Fatalf("word=%d want 1 after space", e.CurrentWord())
	}
	if e.WordError(0) != true {
		t.Fatal("word 0 should be marked wrong")
	}
}

func TestSpaceIgnoredWhenEmpty(t *testing.T) {
	e := New("hello world", DefaultConfig())
	e.Space()
	e.Space()
	if e.CurrentWord() != 0 {
		t.Fatalf("word=%d want 0 (empty space ignored)", e.CurrentWord())
	}
	typeStr(e, "hello")
	e.Space()
	e.Space()
	if e.CurrentWord() != 1 {
		t.Fatalf("word=%d want 1 (double space ignored)", e.CurrentWord())
	}
}

func TestBackspaceStaysAndFixes(t *testing.T) {
	e := New("hello world", DefaultConfig())
	typeStr(e, "hellx")
	e.Backspace()
	typeStr(e, "o")
	if got := e.InputWord(0); got != "hello" {
		t.Fatalf("input=%q want hello", got)
	}
}

func TestBackspaceToPreviousOnlyIfIncorrect(t *testing.T) {
	e := New("hello world", DefaultConfig())
	typeStr(e, "hello")
	e.Space() // correct, locked
	e.Backspace()
	if e.CurrentWord() != 1 {
		t.Fatalf("word=%d want 1 (must not jump back to correct word)", e.CurrentWord())
	}

	e2 := New("hello world", DefaultConfig())
	typeStr(e2, "hellx")
	e2.Space() // incorrect
	e2.Backspace()
	if e2.CurrentWord() != 0 {
		t.Fatalf("word=%d want 0 (can jump back to incorrect word)", e2.CurrentWord())
	}
}

func TestDeleteWordClearsCurrent(t *testing.T) {
	e := New("hello world", DefaultConfig())
	typeStr(e, "hellx")
	e.DeleteWord() // option+backspace / ctrl+w
	if got := e.InputWord(0); got != "" {
		t.Fatalf("input=%q want empty", got)
	}
	if e.CurrentWord() != 0 {
		t.Fatalf("word=%d want 0", e.CurrentWord())
	}
	// Empty current + incorrect previous -> wipes previous.
	typeStr(e, "hellx")
	e.Space()
	e.DeleteWord()
	if e.CurrentWord() != 0 || e.InputWord(0) != "" {
		t.Fatalf("want jump back and wipe, got word=%d input=%q", e.CurrentWord(), e.InputWord(0))
	}
}

func TestAccuracyCountsKeystrokes(t *testing.T) {
	e := New("hello world", DefaultConfig())
	typeStr(e, "hello")
	e.Space()
	if got := e.Accuracy(); got != 1 {
		t.Fatalf("acc=%v want 1", got)
	}
	typeStr(e, "worlx")
	e.Backspace()
	typeStr(e, "d") // fixed, but the wrong keystroke still counts
	if got := e.Accuracy(); got >= 1 {
		t.Fatalf("acc=%v want <1 after corrected mistake", got)
	}
}

func TestWordsDoneOnLastWordNoSpace(t *testing.T) {
	e := New("hi yo", DefaultConfig())
	typeStr(e, "hi")
	e.Space()
	typeStr(e, "yo")
	if !e.Done() {
		t.Fatal("should finish when last word exactly typed")
	}
}
