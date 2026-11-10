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
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						"/opt/fixture/repo",
					),
				),
			),
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
									"main",
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

	t.Run("found issues in real unpushed commits", func(t *testing.T) {
		t.Parallel()

		envs := sharedContainerFixture(t).Env()
		want := []issue.Issue{
			issue.NewInfo(
				`task doesn't match the required pattern '(PROJ|BUG)-[0-9]+'`,
			),
			issue.NewCritical(
				envs.MustGet("LINT_COMMIT_02")[:7] + ": ticket not found in subject by expression '(PROJ|BUG)-[0-9]+'",
			),
			issue.NewCritical(
				envs.MustGet("LINT_COMMIT_01")[:7] + ": ticket not found in subject by expression '(PROJ|BUG)-[0-9]+'",
			),
		}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						"/opt/fixture/repo",
					),
				),
			),
			*config.NewConfig(
				config.WithBranch(
					branch.NewConfig(
						branch.WithTask(
							task.NewConfig(
								task.WithID(
									id.NewConfig(
										value.NewString(
											issue.Info,
											`(PROJ|BUG)-[0-9]+`,
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
									"main",
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
															`(PROJ|BUG)-[0-9]+`,
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

	t.Run("branch name not matches required pattern", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				`branch name 'feature/TASK-123-lint' doesn't match the required pattern '^bugfix/[A-Z]+-\d+'`,
			),
		}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						"/opt/fixture/repo",
					),
				),
			),
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
									`^bugfix/[A-Z]+-\d+`,
								),
							),
						),
						branch.WithDefault(
							defaults.NewConfig(
								defaults.WithName(
									"main",
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

	t.Run("fails on empty config", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						"/opt/fixture/repo",
					),
				),
			),
			*config.NewConfig(),
		).Issues(ctx)
		assert.Nil(t, got)
		assert.EqualError(
			t,
			gotErr,
			"invalid config: invalid branch: invalid default branch: empty default branch name",
		)
	})

	t.Run("fails on canceled context", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(
			context.Background(),
		)
		cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						"/opt/fixture/repo",
					),
				),
			),
			*config.NewConfig(),
		).Issues(ctx)
		assert.Nil(t, got)
		assert.Contains(
			t,
			gotErr.Error(),
			"context canceled",
		)
	})

	t.Run("skips merge commits of an unmerged branch", func(t *testing.T) {
		t.Parallel()
		env := sharedContainerFixture(t).Env()
		want := []issue.Issue{
			issue.NewCritical(
				env.MustGet("RELEASE_BETA_HEAD")[:7] + ": ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
			issue.NewCritical(
				env.MustGet("RELEASE_BETA_01")[:7] + ": ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
			issue.NewCritical(
				env.MustGet("RELEASE_ALPHA_HEAD")[:7] + ": ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
			issue.NewCritical(
				env.MustGet("RELEASE_ALPHA_01")[:7] + ": ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
		}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						env.MustGet("LINT_MERGE_REPO"),
					),
				),
			),
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
									"main",
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
	t.Run("skips revert commit on the current branch", func(t *testing.T) {
		t.Parallel()
		env := sharedContainerFixture(t).Env()
		want := []issue.Issue{}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()
		got, gotErr := impl.NewLinter(
			igit.NewClient(
				git.NewClient(
					shfx.NewShell(
						sharedContainerFixture(t).
							Container(),
						env.MustGet("LINT_REVERT_REPO"),
					),
				),
			),
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
									"main",
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
