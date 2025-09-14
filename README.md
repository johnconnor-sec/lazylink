# lazylink

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/dl/)

A TUI for linking notes in markdown directories. Featuring filtering and search

## Installation

```bash
git clone https://github.com/johnconnor-sec/lazylink.git
cd lazylink
go build -o lazylink ./cmd/lazylink/
```

## Usage

```bash
./lazylink --vault /path/to/markdown/dir/ --ignore-dirs comma,separated,ignored,dirs
```

Or set the `$ZK_NOTEBOOK_DIR` environment variable.

You can place the binary in your `$GOBIN` to make it globally available

Create an alias if you run in the same vault or ignore the same directories repeatedly

```bash
# In your bashrc or zshrc
alias ll="lazylink --vault /path/to/markdown/dir/ --ignore-dirs comma,separated,ignored,dirs"
```

## Features

- Fuzzy search through notes
- Preview target note content
- Insert markdown links into "## Related" section
- Undo last link addition
- Rescan vault for changes
- Filter directories

## Controls

- `↑/k`, `↓/j`: Navigate notes
- `Enter`: Insert link
- `n/p`: Change target note
- `Tab`: Switch focus between panes
- `/`: Enter search mode
- `f`: Filter Directory
- `r`: Rescan vault
- `u`: Undo last link
- `q`: Quit

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)

## License

GPL
