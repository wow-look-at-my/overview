package header

import _ "embed"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummary(t *testing.T) {
	cases := []struct {
		name, path, src, want string
	}{
		{"go package doc", "a.go", "// Package a does one thing. It does it well.\npackage a\n", "Package a does one thing."},
		{"go build tag first", "a.go", "//go:build linux\n\n// Package a runs on linux.\npackage a\n", "Package a runs on linux."},
		{"license block skipped", "a.go", "// Copyright 2026 Someone.\n// SPDX-License-Identifier: MIT\n\n// Package a is real.\npackage a\n", "Package a is real."},
		{"one-line jsdoc", "x.ts", "/** <activity-feed> — a filterable activity/event log table. */\nimport x from 'y';\n", "<activity-feed> — a filterable activity/event log table."},
		{"multi-line block", "x.css", "/* <activity-feed> shadow-DOM styles.\n *\n * Themed through custom properties.\n */\n:host {}\n", "<activity-feed> shadow-DOM styles."},
		{"wrapped sentence joins", "x.js", "// Every container rides the mirror\n// when it starts. Second.\n", "Every container rides the mirror when it starts."},
		{"blank comment line ends paragraph", "x.js", "// First part\n//\n// Second part.\n", "First part"},
		{"abbreviation does not cut", "x.ts", "// Folds e.g. hooks into hook.\n", "Folds e.g. hooks into hook."},
		{"code first means no header", "x.ts", "import x from 'y';\n// late comment\n", ""},
		{"empty file", "x.ts", "", ""},
		{"shebang then hash", "x.sh", testSummarySh, "Installs the tools."},
		{"yaml with no comment", "ci.yml", "name: CI\non: push\n", ""},
		{"yaml comment after directive", "x.yml", "---\n# The deploy config.\nname: x\n", "The deploy config."},
		{"makefile by name", "Makefile", "# Builds it.\nall:\n", "Builds it."},
		{"markdown heading", "README.md", "# overview\n\nText.\n", "overview"},
		{"markdown front matter", "SKILL.md", "---\nname: x\n---\n\n# The skill\n", "The skill"},
		{"markdown with no heading", "notes.md", "Plain notes here. More.\n", "Plain notes here."},
		{"html comment", "index.html", "<!doctype html>\n<!-- The gallery page. -->\n<html>", "The gallery page."},
		{"html title", "index.html", "<!doctype html><html><head><title>Gallery</title>", "Gallery"},
		{"unknown type", "data.json", "{}", ""},
		{"glsl version line first", "x.vert.glsl", "#version 300 es\n// Fullscreen triangle from gl_VertexID.\nout vec2 vUv;\n", "Fullscreen triangle from gl_VertexID."},
		{"glsl extension and precision", "x.frag.glsl", "#version 300 es\n#extension GL_EXT_foo : enable\nprecision highp float;\n/* Tone map. */\n", "Tone map."},
		{"wgsl enable", "x.wgsl", "enable f16;\n// Half-float blur.\n", "Half-float blur."},
		{"crlf", "a.go", "// Package a ends lines with CRLF.\r\npackage a\r\n", "Package a ends lines with CRLF."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, Summary(c.path, []byte(c.src)))
		})
	}
}

func TestCommentable(t *testing.T) {
	assert.True(t, Commentable("src/a.ts"))
	assert.True(t, Commentable("Dockerfile"))
	assert.True(t, Commentable("X.GO"))
	assert.False(t, Commentable("LICENSE"))
	assert.False(t, Commentable("go.sum"))
	assert.False(t, Commentable("logo.png"))
}


//go:embed testdata/test-summary.sh
var testSummarySh string
