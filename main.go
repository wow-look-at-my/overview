// Command overview prints a repository's file tree with each file's header summary.
package main

import (
	"fmt"
	"os"

	"github.com/wow-look-at-my/overview/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "overview:", err)
		os.Exit(1)
	}
}
