package tui

import (
	"strings"

	"github.com/johnconnor-sec/lazylink/internal/notes"
)

// fuzzyMatch performs fuzzy string matching for search functionality
// Returns true if the query matches the text using flexible substring matching
func fuzzyMatch(note notes.Note, query string) bool {
	text := note.Title + " " + note.Content
	if query == "" {
		return true
	}

	text = strings.ToLower(text)
	query = strings.ToLower(query)

	// Use fuzzysearch, it's better
	// RankMatch returns a score (0-1); use max helper
	// threshold := max(len(query)-1, 1)
	// if score := fuzzy.RankMatch(query, text); score > threshold {
	// 	return true
	// }

	// Fallback: Exact match
	if strings.Contains(text, query) {
		return true
	}

	// Fallback: Split query into words and check if all words are present
	queryWords := strings.Fields(query)
	if len(queryWords) > 1 {
		allWordsMatch := true
		for _, word := range queryWords {
			if !strings.Contains(text, word) {
				allWordsMatch = false
				break
			}
		}
		if allWordsMatch {
			return true
		}
	}

	// Fallback: Check for acronym match (first letters of words)
	if len(query) >= 2 {
		textWords := strings.Fields(text)
		if len(textWords) >= len(query) {
			acronym := ""
			for _, word := range textWords {
				if len(word) > 0 {
					acronym += string(word[0])
				}
			}
			if strings.Contains(acronym, query) {
				return true
			}
		}
	}

	return false
}
