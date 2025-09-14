package tui

import "github.com/charmbracelet/glamour"

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
