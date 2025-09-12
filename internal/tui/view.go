package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View controls how the TUI looks
// ==============================================
func (m Model) View() string {
	if m.w == 0 || m.h == 0 {
		return "Loading…"
	}

	title := lipgloss.NewStyle().Bold(true)
	subtle := lipgloss.NewStyle().Faint(true)
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	selected := lipgloss.NewStyle().Bold(true).Underline(true)
	focus := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ok := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errc := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	header := title.Render("Linkr — Obsidian linker")
	help := subtle.Render("[↑/k, ↓/j] move  [Enter] link  [n/p] target  [Tab] focus  [r] rescan  [q] quit")

	// left pane
	var left strings.Builder
	left.WriteString(title.Render("Unlinked Notes") + "\n")
	left.WriteString(subtle.Render(fmt.Sprintf("Target: %s", m.notes[m.targetIdx].Title)) + "\n\n")

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
	leftPane := border.Width(m.w/2 - 2).Height(m.h - 4).Render(left.String())

	// right pane
	var right strings.Builder
	right.WriteString(title.Render("Target Preview") + "\n")
	right.WriteString(subtle.Render(m.notes[m.targetIdx].Path) + "\n\n")
	right.WriteString(m.vp.View())
	rightPane := border.Width(m.w - (m.w / 2) - 2).Height(m.h - 4).Render(right.String())

	status := m.status
	if m.err != nil {
		status = errc.Render(m.err.Error())
	} else if status != "" {
		status = ok.Render(status)
	}

	statusLine := subtle.Render(status)
	row := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	return lipgloss.JoinVertical(lipgloss.Left, header, help, row, statusLine)
}
