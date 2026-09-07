package gitamix

import (
	"github.com/spf13/cobra"
)

var (
	// cmd defines the root command of the gitamix CLI.
	//
	// The root command aggregates the registered subcommands,
	// such as lint, and delegates their execution to the cobra framework.
	cmd = &cobra.Command{
		Use:   "gitamix",
		Short: "CLI tool that improves the user's git using gitamix tools",
		Long: `Gitamix is a CLI tool that improves the user's git
using gitamix tools, such as git commit message linting`,
	}
)

// Execute runs the root command of the gitamix CLI
// and returns an error if the execution fails.
func Execute() error {
	return cmd.Execute()
}
