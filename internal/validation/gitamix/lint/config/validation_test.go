package config_test

import (
	"testing"

	"github.com/gitamix/lint/config"
	branchConfig "github.com/gitamix/lint/config/branch"
	"github.com/gitamix/lint/config/branch/defaults"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/validation/gitamix/lint/config"
)

func TestValidation_Validate(t *testing.T) {
	t.Parallel()

	t.Run("failed on empty default branch name", func(t *testing.T) {
		t.Parallel()
		gotErr := impl.
			NewValidation(
				*config.NewConfig(),
			).
			Validate()
		assert.EqualError(
			t,
			gotErr,
			"empty default branch name",
		)
	})

	t.Run("not failed on complete configuration", func(t *testing.T) {
		t.Parallel()
		gotErr := impl.
			NewValidation(
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
				),
			).
			Validate()
		assert.NoError(t, gotErr)
	})
}
