package lint

import (
	"log"

	"github.com/gitamix/lint/config"
	"github.com/spf13/cobra"
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
			_, err := config.Load(flags.config)
			if err != nil {
				log.Fatalf("ERROR: %s", err)
			}
		},
	}
)

// Command returns the lint command that can be attached
// to the root command of the gitamix CLI.
func Command() *cobra.Command {
	return cmd
}
