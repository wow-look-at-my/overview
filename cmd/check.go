// The check subcommand fails when a source file has no header comment.

package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wow-look-at-my/overview/internal/header"
)

var excludes []string

var checkCmd = &cobra.Command{
	Use:   "check [path...]",
	Short: "Fail when a source file has no header comment",
	Long: `check lists every file whose type has a comment syntax but that carries no
leading summary comment, and exits 1 when it finds one. A generated file
("Code generated ... DO NOT EDIT.") is never flagged.

--exclude takes a glob. It matches a whole path, a base name, or, with a
trailing slash, a directory prefix. Repeat it for more than one.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(c *cobra.Command, args []string) error {
		for _, pat := range excludes {
			if _, err := path.Match(pat, ""); err != nil {
				return fmt.Errorf("--exclude %q: %w", pat, err)
			}
		}
		files, err := readFiles(args)
		if err != nil {
			return err
		}
		checked, missing := 0, 0
		for _, f := range files {
			if !header.Commentable(f.path) || f.generated || excluded(f.path, excludes) {
				continue
			}
			checked++
			if f.summary == "" {
				missing++
				fmt.Fprintf(c.OutOrStdout(), "%s: no header comment\n", f.path)
			}
		}
		if missing > 0 {
			return fmt.Errorf("%d of %d files have no header comment", missing, checked)
		}
		fmt.Fprintf(c.OutOrStdout(), "%d files checked, every one has a header comment\n", checked)
		return nil
	},
}

func init() {
	checkCmd.Flags().StringArrayVar(&excludes, "exclude", nil, "glob of files to skip (repeatable)")
	rootCmd.AddCommand(checkCmd)
}

// excluded reports whether p matches one of the patterns.
func excluded(p string, patterns []string) bool {
	for _, pat := range patterns {
		if strings.HasSuffix(pat, "/") {
			if strings.HasPrefix(p, pat) {
				return true
			}
			continue
		}
		if ok, _ := path.Match(pat, p); ok {
			return true
		}
		if ok, _ := path.Match(pat, path.Base(p)); ok {
			return true
		}
	}
	return false
}
