package tui

import (
	"fmt"
	"linkr/internal/notes"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// recompute the TUI when notes are linked via [enter]
// or when the target changes
func (m *Model) recompute() {
	if len(m.notes) == 0 {
		m.candidates = nil
		return
	}
	target := m.notes[m.targetIdx]
	content := notes.Read(target.Path)

	// Candidates = notes not linked in target AND matching search query
	var cands []notes.Note
	for _, n := range m.notes {
		if n.Path == target.Path {
			continue
		}
		rel := notes.RelPath(filepath.Dir(target.Path), n.Path)
		if !notes.ContainsLink(content, rel, n.Title) {
			// Apply search filter if search mode is active
			if m.searchMode && m.searchQuery != "" {
				if fuzzyMatch(n, m.searchQuery) {
					cands = append(cands, n)
				}
			} else {
				cands = append(cands, n)
			}
		}
	}
	m.candidates = cands
	m.leftIdx = clamp(m.leftIdx, 0, max(0, len(m.candidates)-1))

	// Left pane
	selected := lipgloss.NewStyle().Bold(true).Underline(true)
	focus := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	var left strings.Builder
	left.WriteString(lipgloss.NewStyle().Bold(true).Render("Unlinked Notes") + "\n")
	left.WriteString(lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("Target: %s", m.notes[m.targetIdx].Title)) + "\n\n")

	if len(m.candidates) == 0 {
		left.WriteString("All linked. Nice.\n")
	} else {
		for i, n := range m.candidates {
			item := n.Title
			if i == m.leftIdx && m.focus == focusLeft {
				item = focus.Render("> ") + selected.Render(item)
			} else if i == m.leftIdx {
				item = "> " + item
			} else {
				item = "  " + item
			}
			left.WriteString(item + "\n")
		}
	}
	m.leftVp.SetContent(left.String())
	m.leftVp.SetYOffset(clamp(m.leftIdx-m.leftVp.Height/2, 0, max(0, len(m.candidates)-m.leftVp.Height)))

	// Right pane
	rightPreview := preview(content, 3000)
	m.rightVp.SetContent(rightPreview)
}
