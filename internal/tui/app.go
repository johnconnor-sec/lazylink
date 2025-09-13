// Package tui
package tui

import (
	"linkr/internal/notes"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type focusArea int

const (
	focusLeft focusArea = iota
	focusRight
)

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
}

func New(vault string, all []notes.Note) Model {
	m := Model{
		vault:   vault,
		notes:   all,
		focus:   focusLeft,
		leftVp:  viewport.Model{Width: 60, Height: 40},
		rightVp: viewport.Model{Width: 60, Height: 40},
	}
	m.recompute()
	return m
}

func (m Model) Init() tea.Cmd { return nil }
