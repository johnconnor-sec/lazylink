package tui

import (
	"fmt"
	"path/filepath"

	"github.com/johnconnor-sec/note-linkr/internal/notes"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles the updating the TUI (via recompute)
// and keypresses / commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.leftVp.Width = m.w/2 - 4
		m.rightVp.Width = m.w - (m.w / 2) - 4

		// Adjust viewport height based on search mode
		if m.searchMode {
			m.leftVp.Height = m.h - 11 // Account for search input (3 lines)
			m.rightVp.Height = m.h - 8
		} else {
			m.leftVp.Height = m.h - 8
			m.rightVp.Height = m.h - 8
		}

		m.recompute()
		return m, nil

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.focus == focusLeft {
				m.leftVp.ScrollUp(1)
			} else {
				m.rightVp.ScrollUp(1)
			}
		case tea.MouseWheelDown:
			if m.focus == focusLeft {
				m.leftVp.ScrollDown(1)
			} else {
				m.rightVp.ScrollDown(1)
			}
		}

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.searchMode {
				// Exit search mode
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.Reset()
				m.searchInput.Blur()
				m.leftVp.Height = m.h - 8 // Reset viewport height
				m.status = "Search cleared"
				m.recompute()
			} else {
				return m, tea.Quit
			}
			return m, nil

		case "/":
			if !m.searchMode {
				// Enter search mode
				m.searchMode = true
				m.searchInput.Focus()
				m.leftVp.Height = m.h - 11 // Adjust for search input
				m.status = "Search mode: type to filter notes"
				m.recompute()
			}
			return m, nil

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
			if m.focus == focusLeft {
				if len(m.candidates) > 0 && m.leftIdx > 0 {
					m.leftIdx--
					m.recompute()
				}
			} else {
				m.rightVp.ScrollUp(1)
			}
			return m, nil

		case "down", "j":
			if m.focus == focusLeft {
				if len(m.candidates) > 0 && m.leftIdx < len(m.candidates)-1 {
					m.leftIdx++
					m.recompute()
				}
			} else {
				m.rightVp.ScrollDown(1)
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

	// Handle text input when in search mode
	if m.searchMode {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)

		// Update search query and recompute if changed
		newQuery := m.searchInput.Value()
		if newQuery != m.searchQuery {
			m.searchQuery = newQuery
			m.recompute()
			if m.searchQuery == "" {
				m.status = "Search mode: type to filter notes"
			} else {
				m.status = fmt.Sprintf("Search: %s (%d matches)", m.searchQuery, len(m.candidates))
			}
		}
	}

	m.leftVp, cmd = m.leftVp.Update(msg)
	cmds = append(cmds, cmd)
	m.rightVp, cmd = m.rightVp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
