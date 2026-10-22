package tree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func render(t *testing.T, entries []Entry, depth int) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, Render(&b, entries, depth))
	return b.String()
}

var sample = []Entry{
	{Path: "main.go", Summary: "Command overview prints the tree."},
	{Path: "LICENSE"},
	{Path: "internal/tree/tree.go", Summary: "Package tree draws trees."},
	{Path: "internal/header/header.go", Summary: "Package header reads headers."},
	{Path: "internal/header/header_test.go"},
	{Path: "cmd/root.go", Summary: "Package cmd holds the CLI."},
}

func TestRenderDirsFirstThenFilesSorted(t *testing.T) {
	want := strings.Join([]string{
		"├── cmd/",
		"│   └── root.go  ← Package cmd holds the CLI.",
		"├── internal/",
		"│   ├── header/",
		"│   │   ├── header.go  ← Package header reads headers.",
		"│   │   └── header_test.go",
		"│   └── tree/",
		"│       └── tree.go  ← Package tree draws trees.",
		"├── LICENSE",
		"└── main.go  ← Command overview prints the tree.",
		"",
	}, "\n")
	assert.Equal(t, want, render(t, sample, 0))
}

func TestRenderDepthFoldsWithCount(t *testing.T) {
	want := strings.Join([]string{
		"├── cmd/ (1 files)",
		"├── internal/ (3 files)",
		"├── LICENSE",
		"└── main.go  ← Command overview prints the tree.",
		"",
	}, "\n")
	assert.Equal(t, want, render(t, sample, 1))
}

func TestRenderAlignsSummariesAmongSiblings(t *testing.T) {
	out := render(t, []Entry{
		{Path: "a.go", Summary: "A."},
		{Path: "longer_name.go", Summary: "B."},
	}, 0)
	assert.Equal(t, "├── a.go            ← A.\n└── longer_name.go  ← B.\n", out)
}

func TestRenderCapsPadding(t *testing.T) {
	long := strings.Repeat("x", 60) + ".go"
	out := render(t, []Entry{{Path: "a.go", Summary: "A."}, {Path: long, Summary: "B."}}, 0)
	first, _, _ := strings.Cut(out, "\n")
	assert.Equal(t, "├── a.go"+strings.Repeat(" ", maxPad-len("a.go"))+"  ← A.", first)
	assert.Contains(t, out, long+"  ← B.")
}

func TestFlatSortsFullPaths(t *testing.T) {
	var b strings.Builder
	require.NoError(t, Flat(&b, sample))
	want := strings.Join([]string{
		"LICENSE",
		"cmd/root.go  ← Package cmd holds the CLI.",
		"internal/header/header.go  ← Package header reads headers.",
		"internal/header/header_test.go",
		"internal/tree/tree.go  ← Package tree draws trees.",
		"main.go  ← Command overview prints the tree.",
		"",
	}, "\n")
	assert.Equal(t, want, b.String())
}

func TestRenderEmpty(t *testing.T) {
	assert.Equal(t, "", render(t, nil, 0))
}
