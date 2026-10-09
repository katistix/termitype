package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/katistix/termitype/internal/config"
	"github.com/katistix/termitype/internal/theme"
	"github.com/katistix/termitype/internal/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "theme" {
		if err := runTheme(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "termitype:", err)
			os.Exit(1)
		}
		return
	}

	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Println(config.Usage())
			return
		}
		fmt.Fprintln(os.Stderr, "termitype:", err)
		fmt.Fprintln(os.Stderr, "try: termitype --help")
		os.Exit(1)
	}

	m := ui.New(cfg, theme.Load())
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// runTheme handles `termitype theme [name]`: with a name it sets the theme
// directly, otherwise it opens the fuzzy picker.
func runTheme(args []string) error {
	if len(args) > 0 {
		if args[0] == "-h" || args[0] == "--help" {
			fmt.Println(config.Usage())
			return nil
		}
		name := strings.Join(args, " ")
		t, ok := theme.Get(name)
		if !ok {
			msg := fmt.Sprintf("unknown theme %q", name)
			if ms := theme.Filter(theme.All(), name); len(ms) > 0 {
				var names []string
				for _, m := range ms[:min(5, len(ms))] {
					names = append(names, m.Theme.Name)
				}
				msg += "; did you mean: " + strings.Join(names, ", ")
			}
			return errors.New(msg)
		}
		return saveTheme(t)
	}

	p, err := tea.NewProgram(ui.NewPicker(theme.Load()), tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	t, ok := p.(ui.Picker).Chosen()
	if !ok {
		return nil
	}
	return saveTheme(t)
}

func saveTheme(t theme.Theme) error {
	if err := theme.Save(t.Name); err != nil {
		return err
	}
	fmt.Println("theme set to", t.Display())
	return nil
}
