package lint

// init setups command.
//
// Go invokes it automatically when the package is imported.
func init() {
	cmd.
		PersistentFlags().
		StringVar(
			&flags.config,
			"config",
			"gitamix.lint.yml",
			"path to configuration yaml file",
		)
}
