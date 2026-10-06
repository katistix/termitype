package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/katistix/termitype/internal/config"
	"github.com/katistix/termitype/internal/ui"
)

func main() {
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

	m := ui.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
