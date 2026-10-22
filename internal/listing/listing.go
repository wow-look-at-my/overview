// Package listing writes file paths and their summaries, one full path per line.
package listing

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// Entry is one file and the summary written beside it.
type Entry struct {
	Path    string
	Summary string
}

// Write writes entries sorted by path, each with its summary when it has one.
func Write(w io.Writer, entries []Entry) error {
	sorted := slices.SortedFunc(slices.Values(entries), func(a, b Entry) int { return strings.Compare(a.Path, b.Path) })
	for _, e := range sorted {
		line := e.Path
		if e.Summary != "" {
			line += "  ← " + e.Summary
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
