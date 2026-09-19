//go:build integration
// +build integration

package lint_test

import (
	"context"
	"testing"
	"time"

	"github.com/gitamix/git/client/git"
	"github.com/gitamix/lint/config"
	"github.com/gitamix/lint/config/branch"
	"github.com/gitamix/lint/config/branch/defaults"
	"github.com/gitamix/lint/config/branch/name"
	"github.com/gitamix/lint/config/commit"
	"github.com/gitamix/lint/config/commit/message"
	"github.com/gitamix/lint/config/commit/message/body"
	"github.com/gitamix/lint/config/commit/message/body/mandate"
	"github.com/gitamix/lint/config/commit/message/subject"
	"github.com/gitamix/lint/config/commit/message/subject/description"
	"github.com/gitamix/lint/config/commit/scope"
	"github.com/gitamix/lint/config/commit/types"
	"github.com/gitamix/lint/config/task"
	"github.com/gitamix/lint/config/task/id"
	"github.com/gitamix/lint/config/value"
	"github.com/gitamix/lint/issue"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/app/gitamix/lint"
	igit "github.com/gitamix/cli/internal/client/git"
	shfx "github.com/gitamix/cli/internal/test/fixture/shell"
)

func TestLinter_Issues(t *testing.T) {
	t.Parallel()

	t.Run("clean branch without issues", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{}
		cli := igit.NewClient(
			git.NewClient(
				shfx.NewShell(
					sharedContainerFixture(t).
						Container(),
					"/opt/fixture/repo",
				),
			),
		)
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			cli,
			*config.NewConfig(
				config.WithBranch(
					branch.NewConfig(
						branch.WithTask(
							task.NewConfig(
								task.WithID(
									id.NewConfig(
										value.NewString(
											issue.Info,
											`(TASK|PROJ|BUG)-[0-9]+`,
										),
									),
								),
							),
						),
						branch.WithName(
							name.NewConfig(
								value.NewString(
									issue.Warning,
									`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
								),
							),
						),
						branch.WithDefault(
							defaults.NewConfig(
								defaults.WithName(
									"master",
								),
							),
						),
					),
				),
				config.WithCommit(
					commit.NewConfig(
						commit.WithMessage(
							message.NewConfig(
								message.WithSubject(
									subject.NewConfig(
										subject.WithDescription(
											description.NewConfig(
												description.WithLength(
													value.NewRange(
														issue.Warning,
														10,
														72,
													),
												),
											),
										),
										subject.WithTask(
											task.NewConfig(
												task.WithID(
													id.NewConfig(
														value.NewString(
															issue.Critical,
															`(TASK|PROJ|BUG)-[0-9]+`,
														),
													),
												),
											),
										),
									),
								),
								message.WithBody(
									body.NewConfig(
										body.WithLength(
											value.NewRange(
												issue.Info,
												10,
												255,
											),
										),
										body.WithMandate(
											mandate.NewConfig(
												mandate.WithTypes(
													value.NewStrings(
														issue.Critical,
														"fix",
														"chore",
														"refactor",
														"perf",
													),
												),
											),
										),
									),
								),
							),
						),
						commit.WithScope(
							scope.NewConfig(
								value.NewString(
									issue.Warning,
									`^[A-Za-z _-]+$`,
								),
							),
						),
						commit.WithTypes(
							types.NewConfig(
								value.NewStrings(
									issue.Critical,
									"feat",
									"fix",
									"chore",
									"refactor",
									"perf",
									"test",
									"docs",
								),
							),
						),
					),
				),
			),
		).Issues(ctx)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})
}
