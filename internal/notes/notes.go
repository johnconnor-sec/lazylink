package notes

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Note struct {
	Path    string
	Title   string
	Content string
}

func ScanVault(vault string, ignores []string) ([]Note, error) {
	var out []Note
	err := filepath.WalkDir(vault, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip silently for now
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".obsidian" || strings.HasPrefix(base, ".git") {
				return filepath.SkipDir
			}
			for _, ig := range ignores {
				if base == ig {
					return filepath.SkipDir
				}
			}

			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			rel, _ := filepath.Rel(vault, path)
			content := readContent(path)
			out = append(out, Note{
				Path:    rel,
				Title:   firstTitleOrFilename(path),
				Content: content,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out, nil
}

func firstTitleOrFilename(path string) string {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if after, ok := strings.CutPrefix(line, "# "); ok {
				return strings.TrimSpace(after)
			}
		}
	}
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func Read(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func readContent(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	// Limit to first 10KB to avoid memory issues
	if len(b) > 10240 {
		b = b[:10240]
	}
	return string(b)
}

func RelPath(fromDir, toFile string) string {
	rel, err := filepath.Rel(fromDir, toFile)
	if err != nil {
		return toFile
	}
	return filepath.ToSlash(rel)
}

func ContainsLink(content, rel, title string) bool {
	if strings.Contains(content, "]("+rel+")") {
		return true
	}
	if strings.Contains(content, "[["+title+"]]") {
		return true
	}
	alt := "./" + rel
	return strings.Contains(content, "]("+alt+")")
}

func EnsureLinksSection(content string) string {
	if strings.Contains(content, "## Related") {
		return content
	}
	return strings.TrimRight(content, "\n") + "\n\n## Related\n"
}

func InsertMarkdownLink(path string, linkTitle string, rel string) error {
	content := Read(path)
	content = EnsureLinksSection(content)
	// Avoid dup
	if ContainsLink(content, rel, linkTitle) {
		return nil
	}
	content += "- [" + linkTitle + "](" + rel + ")\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

func RemoveMarkdownLink(path string, linkTitle string, rel string) error {
	content := Read(path)
	if content == "" {
		return nil // File empty or not found, nothing to remove
	}
	lines := strings.Split(content, "\n")
	linkLine := "- [" + linkTitle + "](" + rel + ")"
	var newLines []string
	removed := false
	for _, line := range lines {
		if strings.TrimSpace(line) == linkLine {
			removed = true
			continue // Skip this line
		}
		newLines = append(newLines, line)
	}
	if !removed {
		return nil // Link not found, nothing to remove
	}
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(path, []byte(newContent), 0o644)
}
