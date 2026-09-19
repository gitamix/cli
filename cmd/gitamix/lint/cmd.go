package lint

import (
	"context"
	"fmt"
	"os"

	"github.com/gitamix/git/client/git"
	"github.com/gitamix/lint/config"
	"github.com/sitnikovik/osxec/shell"
	"github.com/spf13/cobra"

	"github.com/gitamix/cli/internal/app/gitamix/lint"
	igit "github.com/gitamix/cli/internal/client/git"
)

var (
	// flags stores the values of the command-line flags of the lint command.
	flags struct {
		// config stores the path to the lint configuration file.
		config string
	}

	// cmd defines the lint command that validates git commit messages.
	cmd = &cobra.Command{
		Use:   "lint",
		Short: "Git commit linter",
		Long: `Git commit linter with support for message format validation,
subject length checks, scope enforcement, and task ID verification`,
		Example: `  gitamix lint --config=/path/to/project/gitamix.lint.yml`,
		Run: func(_ *cobra.Command, _ []string) {
			cfg, err := config.Load(flags.config)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}
			ctx := context.Background()
			cli := igit.NewClient(
				git.NewClient(
					shell.NewShell(),
				),
			)
			issues, err := lint.
				NewLinter(cli, *cfg).
				Issues(ctx)
			if err != nil {
				fmt.Printf(
					"ERROR: %s\n",
					err,
				)
				os.Exit(1)
			}
			fmt.Printf("%d issues\n", len(issues))
			for _, is := range issues {
				fmt.Printf(
					"%s: %s\n",
					is.Type().String(),
					is.Message(),
				)
			}
		},
	}
)

// Command returns the lint command that can be attached
// to the root command of the gitamix CLI.
func Command() *cobra.Command {
	return cmd
}
