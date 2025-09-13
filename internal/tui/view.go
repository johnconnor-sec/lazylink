package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View controls how the TUI looks
func (m Model) View() string {
	if m.w == 0 || m.h == 0 {
		return "Loading…"
	}

	title := lipgloss.NewStyle().Bold(true)
	subtle := lipgloss.NewStyle().Faint(true)
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	ok := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	err := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	header := title.Render("Linkr — Obsidian linker")
	help := subtle.Render("[↑/k, ↓/j] move  [Enter] link  [n/p] target  [Tab] focus  [r] rescan  [q] quit")

	// left pane
	leftPane := border.Width(m.w/2 - 2).Height(m.h - 8).Render(m.leftVp.View())

	// right pane
	var right strings.Builder
	right.WriteString(title.Render("Target Preview") + "\n")
	right.WriteString(subtle.Render(m.notes[m.targetIdx].Path) + "\n\n")
	right.WriteString(m.rightVp.View())
	rightPane := border.Width(m.w - (m.w / 2) - 2).Height(m.h - 8).Render(right.String())

	status := m.status
	if m.err != nil {
		status = err.Render(m.err.Error())
	} else if status != "" {
		status = ok.Render(status)
	}

	statusLine := subtle.Render(status)
	row := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return lipgloss.JoinVertical(lipgloss.Left, header, help, row, statusLine)
}
