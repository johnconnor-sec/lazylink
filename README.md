# note-linkr

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/dl/)

A terminal-based tool for linking notes in markdown directories using fuzzy search and a TUI.

## Installation

Requires Go 1.25+.

```bash
git clone https://github.com/johnconnor-sec/note-linkr.git
cd note-linkr
go build -o linkr ./cmd/linkr
```

## Usage

```bash
./linkr --vault /path/to/markdown/dir/
```

Or set the `ZK_NOTEBOOK_DIR` environment variable.

## Features

- Fuzzy search through notes
- Preview target note content
- Insert markdown links into "## Related" section
- Undo last link addition
- Rescan vault for changes

## Controls

- `↑/k`, `↓/j`: Navigate notes
- `Enter`: Insert link
- `n/p`: Change target note
- `Tab`: Switch focus between panes
- `/`: Enter search mode
- `r`: Rescan vault
- `u`: Undo last link
- `q`: Quit

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)

## License

GPL
