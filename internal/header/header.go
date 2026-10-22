// Package header reads the one-sentence summary from a file's leading comment.
package header

import (
	"path/filepath"
	"regexp"
	"strings"
)

type style int

const (
	styleNone style = iota
	styleSlash
	styleHash
	styleMarkdown
	styleHTML
)

var byExt = map[string]style{
	".go": styleSlash, ".ts": styleSlash, ".tsx": styleSlash, ".mts": styleSlash, ".cts": styleSlash,
	".js": styleSlash, ".jsx": styleSlash, ".mjs": styleSlash, ".cjs": styleSlash,
	".css": styleSlash, ".scss": styleSlash, ".wgsl": styleSlash, ".glsl": styleSlash,
	".c": styleSlash, ".h": styleSlash, ".cc": styleSlash, ".cpp": styleSlash, ".hpp": styleSlash,
	".rs": styleSlash, ".java": styleSlash, ".kt": styleSlash, ".swift": styleSlash, ".proto": styleSlash,
	".sh": styleHash, ".bash": styleHash, ".zsh": styleHash, ".py": styleHash, ".rb": styleHash,
	".yml": styleHash, ".yaml": styleHash, ".toml": styleHash, ".mk": styleHash, ".dats": styleHash,
	".md":   styleMarkdown,
	".html": styleHTML, ".htm": styleHTML,
}

var byName = map[string]style{
	"Makefile": styleHash, "Dockerfile": styleHash, "Containerfile": styleHash,
	"justfile": styleHash, "Justfile": styleHash,
}

func styleOf(path string) style {
	if s, ok := byName[filepath.Base(path)]; ok {
		return s
	}
	return byExt[strings.ToLower(filepath.Ext(path))]
}

// Commentable reports whether this package knows the comment syntax of the file type.
func Commentable(path string) bool { return styleOf(path) != styleNone }


// Summary returns the first sentence of the file's leading comment, or "" when there is none.
func Summary(path string, src []byte) string {
	lines := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")
	var text string
	switch styleOf(path) {
	case styleSlash:
		text = slashSummary(lines)
	case styleHash:
		text = hashSummary(lines)
	case styleMarkdown:
		text = markdownSummary(lines)
	case styleHTML:
		text = htmlSummary(string(src))
	}
	return firstSentence(text)
}

var licenseRe = regexp.MustCompile(`(?i)\b(copyright|spdx-license-identifier|licensed under|all rights reserved)\b`)

// skipDirective reports a line that sits before a header but is not part of it.
func skipDirective(line string) bool {
	for _, p := range []string{"#!", "//go:", "// +build", "// @ts-", "/// <reference", "// eslint-", "// prettier-", "# syntax=", "# shellcheck ", "# -*-", "# vim:", "---", "#version ", "#extension ", "precision ", "enable "} {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

// slashSummary reads the first non-license // or /* */ comment block.
func slashSummary(lines []string) string {
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		switch {
		case line == "" || skipDirective(line):
			i++
		case strings.HasPrefix(line, "//"):
			var block []string
			for ; i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "//"); i++ {
				block = append(block, strings.TrimPrefix(strings.TrimSpace(lines[i]), "//"))
			}
			if text := paragraph(block); text != "" && !licenseRe.MatchString(text) {
				return text
			}
		case strings.HasPrefix(line, "/*"):
			var block []string
			for ; i < len(lines); i++ {
				l := strings.TrimSpace(lines[i])
				end := strings.Contains(l, "*/")
				l, _, _ = strings.Cut(l, "*/")
				l = strings.TrimPrefix(strings.TrimPrefix(l, "/*"), "*")
				block = append(block, l)
				if end {
					i++
					break
				}
			}
			if text := paragraph(block); text != "" && !licenseRe.MatchString(text) {
				return text
			}
		default:
			return ""
		}
	}
	return ""
}

// hashSummary reads the first non-license # comment block.
func hashSummary(lines []string) string {
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		switch {
		case line == "" || skipDirective(line):
			i++
		case strings.HasPrefix(line, "#"):
			var block []string
			for ; i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "#") && !skipDirective(strings.TrimSpace(lines[i])); i++ {
				block = append(block, strings.TrimLeft(strings.TrimSpace(lines[i]), "#"))
			}
			if text := paragraph(block); text != "" && !licenseRe.MatchString(text) {
				return text
			}
		default:
			return ""
		}
	}
	return ""
}

var headingRe = regexp.MustCompile(`^#{1,6}\s+(.*?)\s*#*$`)

// markdownSummary reads the first heading, after any front matter.
func markdownSummary(lines []string) string {
	i := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i = 1; i < len(lines) && strings.TrimSpace(lines[i]) != "---"; i++ {
		}
		i++
	}
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "<!--") {
			continue
		}
		if m := headingRe.FindStringSubmatch(line); m != nil {
			return m[1]
		}
		var block []string
		for ; i < len(lines) && strings.TrimSpace(lines[i]) != ""; i++ {
			block = append(block, lines[i])
		}
		return paragraph(block)
	}
	return ""
}

var (
	htmlCommentRe = regexp.MustCompile(`(?s)^\s*(?:<!doctype[^>]*>\s*)?<!--(.*?)-->`)
	titleRe       = regexp.MustCompile(`(?is)<title>(.*?)</title>`)
)

// htmlSummary reads a leading <!-- --> comment, else the <title>.
func htmlSummary(src string) string {
	if m := htmlCommentRe.FindStringSubmatch(src); m != nil {
		return paragraph(strings.Split(m[1], "\n"))
	}
	if m := titleRe.FindStringSubmatch(src); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// paragraph joins the lines of a comment up to its first blank line.
func paragraph(lines []string) string {
	var words []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			if len(words) > 0 {
				break
			}
			continue
		}
		words = append(words, strings.Fields(l)...)
	}
	return strings.Join(words, " ")
}

var sentenceEndRe = regexp.MustCompile(`[.!?](\s+[A-Z(]|$)`)

// firstSentence cuts text after its first sentence end.
func firstSentence(text string) string {
	if loc := sentenceEndRe.FindStringIndex(text); loc != nil {
		return text[:loc[0]+1]
	}
	return text
}
