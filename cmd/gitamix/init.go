package gitamix

import "github.com/gitamix/cli/cmd/gitamix/lint"

// init setups command.
//
// Go invokes it automatically when the package is imported.
func init() {
	cmd.AddCommand(lint.Command())
}
