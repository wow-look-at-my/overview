// Package cmd holds the overview command line.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/wow-look-at-my/overview/internal/tree"
)

var depth int

var rootCmd = &cobra.Command{
	Use:   "overview [path...]",
	Short: "Print the repository tree with each file's header summary",
	Long: `overview prints the files git knows about as a tree. Beside each file it
prints the first sentence of that file's leading comment. Run it from
inside a git work tree. Paths limit the tree, as they limit git ls-files.`,
	// Without this, cobra reads a path such as "src" as an unknown subcommand.
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(c *cobra.Command, args []string) error {
		files, err := readFiles(args)
		if err != nil {
			return err
		}
		entries := make([]tree.Entry, len(files))
		for i, f := range files {
			entries[i] = tree.Entry{Path: f.path, Summary: f.summary}
		}
		return tree.Render(c.OutOrStdout(), entries, depth)
	},
}

func init() {
	rootCmd.Flags().IntVar(&depth, "depth", 0, "fold directories below this depth into a file count (0 = no limit)")
}

// Execute runs the command line.
func Execute() error { return rootCmd.Execute() }
