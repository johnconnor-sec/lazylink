package tui

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/johnconnor-sec/lazylink/internal/notes"
)

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func preview(s string, maxChars int) string {
	rendered, err := glamour.Render(s, "dark")
	if err != nil {
		rendered = s // Fallback to plaintext if err
	}
	if len(s) <= maxChars {
		return rendered
	}
	return s[:maxChars] + "\n…"
}

func (m *Model) linkSelectedNote() {
	target := m.notes[m.targetIdx]
	selected := m.candidates[m.leftIdx]
	rel := notes.RelPath(filepath.Dir(target.Path), selected.Path)
	if err := notes.InsertMarkdownLink(filepath.Join(m.vault, target.Path), selected.Title, rel); err != nil {
		m.err = err
		m.status = "Insert failed"
	} else {
		m.status = fmt.Sprintf("Linked: %s → %s", target.Title, selected.Title)
		// Push to undo stack
		m.undoStack = append(m.undoStack, UndoAction{TargetPath: target.Path, LinkTitle: selected.Title, Rel: rel})
		m.recompute()
	}
}
