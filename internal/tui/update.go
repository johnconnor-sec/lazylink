package tui

import (
	"fmt"
	"path/filepath"

	"github.com/johnconnor-sec/lazylink/internal/notes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
				m.searchInput.Blur()
				m.leftVp.Height = m.h - 8 // Reset viewport height
				m.status = "Search mode exited"
				// Don't clear searchQuery, but recompute to preserve filtered results
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
				m.status = lipgloss.NewStyle().Render("Search mode: type to filter notes")
				m.recompute()
			} else {
				// Exit search mode and clear results
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.Reset()
				m.searchInput.Blur()
				m.leftVp.Height = m.h - 8 // Reset viewport height
				m.status = "Search cleared"
				m.recompute()
			}
			return m, nil

		case "tab":
			if !m.searchMode {
				if m.focus == focusLeft {
					m.focus = focusRight
				} else {
					m.focus = focusLeft
				}
				m.status = fmt.Sprintf("Focus: %v", m.focus)
				return m, nil
			}

		case "r":
			if !m.searchMode {
				var ignores []string // For now, empty; could store in model if needed
				if all, err := notes.ScanVault(m.vault, ignores); err == nil {
					m.notes = all
					m.targetIdx = 0
					m.leftIdx = 0
					m.filterDir = "" // Clear filter on rescan
					m.recompute()
					m.status = "Rescanned."
				}
				return m, nil
			}

		case "f":
			if !m.searchMode {
				if m.filterDir == "" {
					// Set filter to current target's dir
					m.filterDir = filepath.Dir(m.notes[m.targetIdx].Path)
					m.status = fmt.Sprintf("Filtered to directory: %s", m.filterDir)
				} else {
					// Clear filter
					m.filterDir = ""
					m.status = "Filter cleared"
					m.targetIdx = 0 // Only reset when clearing
				}
				m.leftIdx = 0
				m.recompute()
				return m, nil
			}

		case "n":
			if !m.searchMode {
				// Cycle through filtered targets if filter active
				var targets []int
				for i, n := range m.notes {
					if m.filterDir == "" || filepath.Dir(n.Path) == m.filterDir {
						targets = append(targets, i)
					}
				}
				if len(targets) == 0 {
					return m, nil
				}
				currentTargetIdx := -1
				for j, idx := range targets {
					if idx == m.targetIdx {
						currentTargetIdx = j
						break
					}
				}
				if currentTargetIdx == -1 {
					currentTargetIdx = 0
				}
				nextIdx := (currentTargetIdx + 1) % len(targets)
				m.targetIdx = targets[nextIdx]
				m.leftIdx = 0
				m.recompute()
				return m, nil
			}

		case "p":
			if !m.searchMode {
				// Cycle through filtered targets if filter active
				var targets []int
				for i, n := range m.notes {
					if m.filterDir == "" || filepath.Dir(n.Path) == m.filterDir {
						targets = append(targets, i)
					}
				}
				if len(targets) == 0 {
					return m, nil
				}
				currentTargetIdx := -1
				for j, idx := range targets {
					if idx == m.targetIdx {
						currentTargetIdx = j
						break
					}
				}
				if currentTargetIdx == -1 {
					currentTargetIdx = 0
				}
				prevIdx := (currentTargetIdx - 1 + len(targets)) % len(targets)
				m.targetIdx = targets[prevIdx]
				m.leftIdx = 0
				m.recompute()
				return m, nil
			}

		case "up", "k":
			if !m.searchMode {
				if m.focus == focusLeft {
					if len(m.candidates) > 0 && m.leftIdx > 0 {
						m.leftIdx--
						m.recompute()
						return m, nil
					}
				} else {
					m.rightVp.ScrollUp(1)
				}
				return m, nil
			}

		case "down", "j":
			if !m.searchMode {
				if m.focus == focusLeft {
					if len(m.candidates) > 0 && m.leftIdx < len(m.candidates)-1 {
						m.leftIdx++
						m.recompute()
						return m, nil
					}
				} else {
					m.rightVp.ScrollDown(1)
				}
				return m, nil
			}

		case "u":
			if !m.searchMode && len(m.undoStack) > 0 {
				last := m.undoStack[len(m.undoStack)-1]
				m.undoStack = m.undoStack[:len(m.undoStack)-1] // Pop
				if err := notes.RemoveMarkdownLink(filepath.Join(m.vault, last.TargetPath), last.LinkTitle, last.Rel); err != nil {
					m.err = err
					m.status = "Undo failed"
				} else {
					m.status = fmt.Sprintf("Undid link: %s", last.LinkTitle)
					m.recompute()
				}
				return m, nil
			} else if !m.searchMode {
				m.status = "Nothing to undo"
				return m, nil
			}

		case "enter":
			if !m.searchMode && m.focus == focusLeft && len(m.candidates) > 0 {
				m.linkSelectedNote()
				return m, nil
			}
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
