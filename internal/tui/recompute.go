package tui

import (
	"linkr/internal/notes"
	"path/filepath"
)

// recompute the TUI when notes are linked via [enter]
// or when the target changes
// ==============================================
func (m *Model) recompute() {
	if len(m.notes) == 0 {
		m.candidates = nil
		m.rightPreview = ""
		return
	}
	target := m.notes[m.targetIdx]
	content := notes.Read(target.Path)

	// Candidates = notes not linked in target
	var cands []notes.Note
	for _, n := range m.notes {
		if n.Path == target.Path {
			continue
		}
		rel := notes.RelPath(filepath.Dir(target.Path), n.Path)
		if !notes.ContainsLink(content, rel, n.Title) {
			cands = append(cands, n)
		}
	}
	m.candidates = cands
	m.leftIdx = clamp(m.leftIdx, 0, max(0, len(m.candidates)-1))
	m.rightPreview = preview(content, 3000)
	m.vp.SetContent(m.rightPreview)
}
