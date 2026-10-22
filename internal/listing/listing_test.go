package listing

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteSortsFullPaths(t *testing.T) {
	var b strings.Builder
	require.NoError(t, Write(&b, []Entry{
		{Path: "main.go", Summary: "Command overview prints the list."},
		{Path: "LICENSE"},
		{Path: "internal/header/header_test.go"},
		{Path: "internal/header/header.go", Summary: "Package header reads headers."},
		{Path: "cmd/root.go", Summary: "Package cmd holds the CLI."},
	}))
	want := strings.Join([]string{
		"LICENSE",
		"cmd/root.go  ← Package cmd holds the CLI.",
		"internal/header/header.go  ← Package header reads headers.",
		"internal/header/header_test.go",
		"main.go  ← Command overview prints the list.",
		"",
	}, "\n")
	assert.Equal(t, want, b.String())
}

func TestWriteEmpty(t *testing.T) {
	var b strings.Builder
	require.NoError(t, Write(&b, nil))
	assert.Equal(t, "", b.String())
}
