package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnconnor-sec/lazylink/internal/notes"
	"github.com/johnconnor-sec/lazylink/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	vaultFlag := flag.String("vault", "", "path to Obsidian vault (or set ZK_NOTEBOOK_DIR env)")
	ignoreFlag := flag.String("ignore-dir", "", "comma-separated list of directories to ignore (e.g., .obsidian,.git)")
	flag.Parse()

	vault := *vaultFlag
	if vault == "" {
		vault = os.Getenv("ZK_NOTEBOOK_DIR")
	}
	if vault == "" {
		wd, _ := os.Getwd()
		vault = wd
	}

	vault = filepath.Clean(vault)
	info, err := os.Stat(vault)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Vault not found or not a directory: %s\n", vault)
		os.Exit(1)
	}

	// Parse ignore list
	var ignores []string
	if *ignoreFlag != "" {
		ignores = strings.Split(*ignoreFlag, ",")
		for i, ig := range ignores {
			ignores[i] = strings.TrimSpace(ig)
		}
	}

	// Load notes once at start.
	all, err := notes.ScanVault(vault, ignores)
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(1)
	}
	if len(all) == 0 {
		fmt.Fprintln(os.Stderr, "No markdown notes found.")
		os.Exit(0)
	}

	m := tui.New(vault, all)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
