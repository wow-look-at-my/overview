// Package cmd holds the overview command line.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/wow-look-at-my/overview/internal/listing"
)

var rootCmd = &cobra.Command{
	Use:   "overview [path...]",
	Short: "List the repository's files with each file's header summary",
	Long: `overview lists the files git knows about, one full path per line. Beside
each file it prints the first sentence of that file's leading comment. Run
it from inside a git work tree. Paths limit the list, as they limit git
ls-files.`,
	// Without this, cobra reads a path such as "src" as an unknown subcommand.
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(c *cobra.Command, args []string) error {
		files, err := readFiles(args)
		if err != nil {
			return err
		}
		entries := make([]listing.Entry, len(files))
		for i, f := range files {
			entries[i] = listing.Entry{Path: f.path, Summary: f.summary}
		}
		return listing.Write(c.OutOrStdout(), entries)
	},
}

// Execute runs the command line.
func Execute() error { return rootCmd.Execute() }
