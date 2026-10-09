// Package config parses command-line settings for termitype.
//
// Examples:
//
//	termitype                  # time mode, 15s (default)
//	termitype --mode time --time 30
//	termitype --mode words --words 50
//	termitype -w 25            # shorthand, implies words mode
//	termitype -t 60            # shorthand, implies time mode
package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Mode selects the test type.
type Mode string

const (
	ModeTime  Mode = "time"
	ModeWords Mode = "words"
)

// Config holds all user-facing settings.
type Config struct {
	Mode    Mode // time or words
	Seconds int  // test length in time mode
	Words   int  // test length in words mode
}

const (
	DefaultSeconds = 15
	DefaultWords   = 30
)

// Default returns the default config: 15s time mode.
func Default() Config {
	return Config{Mode: ModeTime, Seconds: DefaultSeconds, Words: DefaultWords}
}

// Usage returns help text.
func Usage() string {
	return `termitype - monkeytype in your terminal

usage:
  termitype [-t seconds | -w count]
  termitype theme [name]

options:
  -t, --time <s>    timed test, s seconds (default 15)
  -w, --words <n>   word-count test, n words (default 30)
  -h, --help        show this help

keys:
  esc               restart test
  ctrl+c            quit
  ctrl+w            delete word (also opt+backspace)

examples:
  termitype           15 second test
  termitype -t 30     30 second test
  termitype -w 50     50 word test
  termitype theme     pick a theme (fuzzy search, live preview)
  termitype theme nord  set the theme directly`
}

// Parse parses args (without program name) into a Config.
func Parse(args []string) (Config, error) {
	cfg := Default()
	modeStr := string(cfg.Mode)

	fs := flag.NewFlagSet("termitype", flag.ContinueOnError)
	fs.StringVar(&modeStr, "mode", modeStr, "test mode: time|words")
	fs.IntVar(&cfg.Seconds, "time", cfg.Seconds, "time limit in seconds")
	fs.IntVar(&cfg.Words, "words", cfg.Words, "word count")
	// Shorthands sharing the same backing vars.
	fs.IntVar(&cfg.Seconds, "t", cfg.Seconds, "shorthand for --time")
	fs.IntVar(&cfg.Words, "w", cfg.Words, "shorthand for --words")

	// Suppress default usage output; caller prints Usage().
	fs.Usage = func() {}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return cfg, err
		}
		return Config{}, fmt.Errorf("%w", err)
	}
	if fs.NArg() > 0 {
		return Config{}, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}

	modeStr = strings.ToLower(strings.TrimSpace(modeStr))
	modeFlagSet := hasFlag(args, "mode")
	timeFlagSet := hasFlag(args, "time", "t")
	wordsFlagSet := hasFlag(args, "words", "w")

	// Shorthands imply their mode when --mode was not given explicitly.
	if !modeFlagSet {
		switch {
		case wordsFlagSet && !timeFlagSet:
			modeStr = string(ModeWords)
		case timeFlagSet && !wordsFlagSet:
			modeStr = string(ModeTime)
		}
	}

	cfg.Mode = Mode(modeStr)
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate checks ranges.
func (c Config) Validate() error {
	switch c.Mode {
	case ModeTime, ModeWords:
	default:
		return fmt.Errorf("invalid mode %q: must be %q or %q", c.Mode, ModeTime, ModeWords)
	}
	if c.Seconds < 5 || c.Seconds > 300 {
		return fmt.Errorf("invalid time %d: must be 5..300 seconds", c.Seconds)
	}
	if c.Words < 5 || c.Words > 500 {
		return fmt.Errorf("invalid words %d: must be 5..500", c.Words)
	}
	return nil
}

// Label returns a short "time 15s" / "words 30" description.
func (c Config) Label() string {
	if c.Mode == ModeWords {
		return fmt.Sprintf("words %d", c.Words)
	}
	return fmt.Sprintf("time %ds", c.Seconds)
}

func hasFlag(args []string, names ...string) bool {
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
		set["-"+n] = true
	}
	for _, a := range args {
		// Split --flag=value.
		name := strings.SplitN(strings.TrimLeft(a, "-"), "=", 2)[0]
		if !strings.HasPrefix(a, "-") {
			continue
		}
		// Match -t / --time style; single-dash long flags also accepted.
		if set[name] || set[strings.TrimLeft(a, "-")] {
			return true
		}
	}
	return false
}
