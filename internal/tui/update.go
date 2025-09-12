package tui

import (
	"fmt"
	"linkr/internal/notes"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles the updating the TUI (via recompute)
// and keypresses / commands
// ==============================================
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.vp.Width = m.w/2 - 4
		m.vp.Height = m.h - 6
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "ESC":
			return m, tea.Quit
		case "tab":
			if m.focus == focusLeft {
				m.focus = focusRight
			} else {
				m.focus = focusLeft
			}
			m.status = fmt.Sprintf("Focus: %v", m.focus)
			return m, nil
		case "r":
			if all, err := notes.ScanVault(m.vault); err == nil {
				m.notes = all
				m.targetIdx = 0
				m.leftIdx = 0
				m.recompute()
				m.status = "Rescanned."
			}
			return m, nil
		case "n":
			if len(m.notes) > 0 {
				m.targetIdx = (m.targetIdx + 1) % len(m.notes)
				m.leftIdx = 0
				m.recompute()
			}
			return m, nil
		case "p":
			if len(m.notes) > 0 {
				m.targetIdx = (m.targetIdx - 1 + len(m.notes)) % len(m.notes)
				m.leftIdx = 0
				m.recompute()
			}
			return m, nil
		case "up", "k":
			if m.focus == focusLeft && len(m.candidates) > 0 && m.leftIdx > 0 {
				m.leftIdx--
			}
			return m, nil
		case "down", "j":
			if m.focus == focusLeft && len(m.candidates) > 0 && m.leftIdx < len(m.candidates)-1 {
				m.leftIdx++
			}
			return m, nil
		case "enter":
			if m.focus == focusLeft && len(m.candidates) > 0 {
				target := m.notes[m.targetIdx]
				selected := m.candidates[m.leftIdx]
				rel := notes.RelPath(filepath.Dir(target.Path), selected.Path)
				if err := notes.InsertMarkdownLink(target.Path, selected.Title, rel); err != nil {
					m.err = err
					m.status = "Insert failed"
				} else {
					m.status = fmt.Sprintf("Linked: %s → %s", target.Title, selected.Title)
					m.recompute()
				}
			}
			return m, nil
		}
	}
	return m, nil
}
