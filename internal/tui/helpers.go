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
		// Invalidate cache for the modified file
		notes.InvalidateCache(m.fileCache, filepath.Join(m.vault, target.Path))
		m.recompute()
	}
}

func (m *Model) cycleTarget(delta int) {
	var targets []int
	for i, n := range m.notes {
		if m.filterDir == "" || filepath.Dir(n.Path) == m.filterDir {
			targets = append(targets, i)
		}
	}
	if len(targets) == 0 {
		return
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
	nextIdx := (currentTargetIdx + delta + len(targets)) % len(targets)
	m.targetIdx = targets[nextIdx]
	m.leftIdx = 0
	m.recompute()
}
