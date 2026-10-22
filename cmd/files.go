package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/wow-look-at-my/overview/internal/header"
)

// readLimit is how much of each file is read. A header sits at the top.
const readLimit = 16 << 10

type file struct {
	path      string
	summary   string
	generated bool
}

// listFiles lists the tracked and untracked-but-not-ignored files under paths, as git sees them.
func listFiles(paths []string) ([]string, error) {
	args := append([]string{"ls-files", "-z", "--cached", "--others", "--exclude-standard", "--"}, paths...)
	var stderr bytes.Buffer
	c := exec.Command("git", args...)
	c.Stderr = &stderr
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	seen := map[string]bool{}
	var files []string
	for _, p := range strings.Split(string(out), "\x00") {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		// A file deleted from the work tree but still in the index has no content to read.
		if st, err := os.Stat(p); err != nil || !st.Mode().IsRegular() {
			continue
		}
		files = append(files, p)
	}
	return files, nil
}

// readFiles lists the files under paths and reads each's header.
func readFiles(paths []string) ([]file, error) {
	names, err := listFiles(paths)
	if err != nil {
		return nil, err
	}
	out := make([]file, 0, len(names))
	for _, name := range names {
		f := file{path: name}
		if header.Commentable(name) {
			src, err := readHead(name)
			if err != nil {
				return nil, err
			}
			f.summary = header.Summary(name, src)
			f.generated = header.Generated(src)
		}
		out = append(out, f)
	}
	return out, nil
}

func readHead(name string) ([]byte, error) {
	fh, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return io.ReadAll(io.LimitReader(fh, readLimit))
}
