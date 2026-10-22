// Package tree draws file paths and their summaries as a box-drawing tree.
package tree

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

// Entry is one file and the summary drawn beside it.
type Entry struct {
	Path    string
	Summary string
}

// maxPad caps the name column, so one long name does not push every summary right.
const maxPad = 40

type node struct {
	name     string
	summary  string
	children map[string]*node
	files    int
}

func (n *node) isDir() bool { return n.children != nil }

// Render writes entries as a tree. A depth above no folds deeper directories into a file count.
func Render(w io.Writer, entries []Entry, depth int) error {
	root := &node{children: map[string]*node{}}
	for _, e := range entries {
		parts := strings.Split(strings.Trim(e.Path, "/"), "/")
		cur := root
		for i, p := range parts {
			cur.files++
			next, ok := cur.children[p]
			if !ok {
				next = &node{name: p}
				if i < len(parts)-1 {
					next.children = map[string]*node{}
				}
				cur.children[p] = next
			}
			cur = next
		}
		cur.summary = e.Summary
	}
	return renderChildren(w, root, "", 1, depth)
}

func sorted(n *node) []*node {
	out := make([]*node, 0, len(n.children))
	for _, c := range n.children {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].isDir() != out[j].isDir() {
			return out[i].isDir()
		}
		return out[i].name < out[j].name
	})
	return out
}

func label(n *node, depth, maxDepth int) string {
	if !n.isDir() {
		return n.name
	}
	if maxDepth > 0 && depth >= maxDepth {
		return fmt.Sprintf("%s/ (%d files)", n.name, n.files)
	}
	return n.name + "/"
}

func renderChildren(w io.Writer, n *node, prefix string, depth, maxDepth int) error {
	kids := sorted(n)
	pad := 0
	for _, c := range kids {
		if c.summary != "" {
			pad = max(pad, utf8.RuneCountInString(label(c, depth, maxDepth)))
		}
	}
	pad = min(pad, maxPad)
	for i, c := range kids {
		branch, indent := "├── ", "│   "
		if i == len(kids)-1 {
			branch, indent = "└── ", "    "
		}
		line := prefix + branch + label(c, depth, maxDepth)
		if c.summary != "" {
			gap := max(pad-utf8.RuneCountInString(label(c, depth, maxDepth)), 0)
			line += strings.Repeat(" ", gap) + "  ← " + c.summary
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
		if c.isDir() && (maxDepth == 0 || depth < maxDepth) {
			if err := renderChildren(w, c, prefix+indent, depth+1, maxDepth); err != nil {
				return err
			}
		}
	}
	return nil
}
