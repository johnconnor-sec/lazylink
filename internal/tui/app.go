// Package tui
package tui

import (
	"github.com/johnconnor-sec/lazylink/internal/notes"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type focusArea int

const (
	focusLeft focusArea = iota
	focusRight
)

type UndoAction struct {
	TargetPath string
	LinkTitle  string
	Rel        string
}

type Model struct {
	vault      string
	notes      []notes.Note
	targetIdx  int
	candidates []notes.Note
	leftIdx    int
	focus      focusArea
	status     string
	err        error

	w, h    int
	leftVp  viewport.Model
	rightVp viewport.Model

	// Search functionality
	searchInput textinput.Model
	searchMode  bool
	searchQuery string

	// Undo functionality
	undoStack []UndoAction
}

func New(vault string, all []notes.Note) Model {
	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search notes..."
	ti.CharLimit = 50
	ti.Width = 30

	m := Model{
		vault:       vault,
		notes:       all,
		focus:       focusLeft,
		leftVp:      viewport.Model{Width: 60, Height: 40},
		rightVp:     viewport.Model{Width: 60, Height: 40},
		searchInput: ti,
		searchMode:  false,
		searchQuery: "",
	}
	m.recompute()
	return m
}

func (m Model) Init() tea.Cmd { return nil }
