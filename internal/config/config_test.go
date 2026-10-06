package config

import (
	"errors"
	"flag"
	"testing"
)

func TestDefaultIs15s(t *testing.T) {
	cfg, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeTime || cfg.Seconds != 15 {
		t.Fatalf("got %+v, want time 15s", cfg)
	}
}

func TestWordsMode(t *testing.T) {
	cfg, err := Parse([]string{"--mode", "words", "--words", "50"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeWords || cfg.Words != 50 {
		t.Fatalf("got %+v", cfg)
	}
}

func TestShorthandImpliesMode(t *testing.T) {
	cfg, err := Parse([]string{"-w", "25"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeWords {
		t.Fatalf("got %+v, want words mode", cfg)
	}
	cfg, err = Parse([]string{"-t", "60"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeTime || cfg.Seconds != 60 {
		t.Fatalf("got %+v, want time 60", cfg)
	}
}

func TestInvalid(t *testing.T) {
	if _, err := Parse([]string{"--mode", "bogus"}); err == nil {
		t.Fatal("want error for bad mode")
	}
	if _, err := Parse([]string{"--time", "1"}); err == nil {
		t.Fatal("want error for tiny time")
	}
	if _, err := Parse([]string{"extra"}); err == nil {
		t.Fatal("want error for positional arg")
	}
}

func TestHelpPassthrough(t *testing.T) {
	if _, err := Parse([]string{"--help"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("want ErrHelp, got %v", err)
	}
}
