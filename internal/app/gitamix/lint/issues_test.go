package lint_test

import (
	"context"
	"testing"

	"github.com/gitamix/lint/config"
	branchConfig "github.com/gitamix/lint/config/branch"
	"github.com/gitamix/lint/config/branch/defaults"
	"github.com/gitamix/lint/config/branch/name"
	commitConfig "github.com/gitamix/lint/config/commit"
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
	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/app/gitamix/lint"
	mockGit "github.com/gitamix/cli/internal/test/mock/client/git"
)

func TestLinter_Issues(t *testing.T) {
	t.Parallel()

	t.Run("found branch not matches required pattern", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				`branch name 'TASK-123' doesn't match the required pattern '^\d+[A-Za-z_-]+$'`,
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"fe85083dd422a6ee60443994f1e51d8fbbac0001",
									commit.ParseMessage(
										[]byte("[TASK-123] feat(ui): add new feature\n\nbody text"),
									),
								),
							}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^\d+[A-Za-z_-]+$`,
									),
								),
							),
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
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("not found on empty commits in branch", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
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
													0,
													255,
												),
											),
										),
									),
								),
							),
							commitConfig.WithScope(
								scope.NewConfig(
									value.NewString(
										issue.Warning,
										`^[A-Za-z _-]+$`,
									),
								),
							),
							commitConfig.WithTypes(
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit without required desc", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewCritical("3dd0e0b: body is required for type 'fix'"),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"fe85083dd422a6ee60443994f1e51d8fbbac0001",
									commit.ParseMessage(
										[]byte("[TASK-123] feat(ui): add new feature\n\nbody text"),
									),
								),
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-123] fix(domain): payment state machine"),
									),
								),
							}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
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
													0,
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
							commitConfig.WithScope(
								scope.NewConfig(
									value.NewString(
										issue.Warning,
										`^[A-Za-z _-]+$`,
									),
								),
							),
							commitConfig.WithTypes(
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("not found when commit task differs the subject one", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-546] feat(ui): add new payment button"),
									),
								),
							}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
								message.NewConfig(
									message.WithSubject(
										subject.NewConfig(
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
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("failed to get current branch", func(t *testing.T) {
		t.Parallel()
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.Branch{}, assert.AnError
						},
					),
				),
				*config.NewConfig(),
			).
			Issues(context.Background())
		assert.Nil(t, got)
		assert.EqualError(
			t,
			gotErr,
			"failed to get current branch: assert.AnError general error for testing",
		)
	})

	t.Run("failed to get merge base", func(t *testing.T) {
		t.Parallel()
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash(""), assert.AnError
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^\d+[A-Za-z_-]+$`,
									),
								),
							),
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
			Issues(context.Background())
		assert.Nil(t, got)
		assert.EqualError(
			t,
			gotErr,
			"failed to get merge base: assert.AnError general error for testing",
		)
	})

	t.Run("failed to get commits", func(t *testing.T) {
		t.Parallel()
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return nil, assert.AnError
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^\d+[A-Za-z_-]+$`,
									),
								),
							),
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
			Issues(context.Background())
		assert.Nil(t, got)
		assert.EqualError(
			t,
			gotErr,
			"failed to get commits: assert.AnError general error for testing",
		)
	})

	t.Run("not found on valid branch and commit", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"fe85083dd422a6ee60443994f1e51d8fbbac0001",
									commit.ParseMessage(
										[]byte("[TASK-123] feat(ui): add new feature\n\nbody text"),
									),
								),
							}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
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
													0,
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
							commitConfig.WithScope(
								scope.NewConfig(
									value.NewString(
										issue.Warning,
										`^[A-Za-z _-]+$`,
									),
								),
							),
							commitConfig.WithTypes(
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found branch task not matches required pattern", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewInfo(
				`task doesn't match the required pattern '(TASK|PROJ|BUG)-[0-9]+'`,
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/user-login"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^feature/[a-z-]+$`,
									),
								),
							),
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
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found branch name and task not match required patterns", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				`branch name 'bad_branch' doesn't match the required pattern '^(feature|bugfix|hotfix)/[A-Z]+-\d+'`,
			),
			issue.NewInfo(
				`task doesn't match the required pattern '(TASK|PROJ|BUG)-[0-9]+'`,
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("bad_branch"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
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
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found multiple issues in one commit in order", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewCritical(
				"3dd0e0b: ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
			issue.NewCritical(
				"3dd0e0b: type must be one of [feat,fix,chore,refactor,perf,test,docs]",
			),
			issue.NewWarning(
				"3dd0e0b: scope not found in subject by expression '^[A-Za-z _-]+$'",
			),
			issue.NewWarning(
				"3dd0e0b: subject description length is not in range [10-72]",
			),
			issue.NewInfo(
				"3dd0e0b: body is not in range [0-5]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("wrong(123): hi\n\nbody text"),
									),
								),
							}, nil
						},
					),
				),
				*config.NewConfig(
					config.WithBranch(
						branchConfig.NewConfig(
							branchConfig.WithTask(
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
							branchConfig.WithName(
								name.NewConfig(
									value.NewString(
										issue.Warning,
										`^(feature|bugfix|hotfix)/[A-Z]+-\d+`,
									),
								),
							),
							branchConfig.WithDefault(
								defaults.NewConfig(
									defaults.WithName(
										"master",
									),
								),
							),
						),
					),
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
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
													0,
													5,
												),
											),
										),
									),
								),
							),
							commitConfig.WithScope(
								scope.NewConfig(
									value.NewString(
										issue.Warning,
										`^[A-Za-z _-]+$`,
									),
								),
							),
							commitConfig.WithTypes(
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found issues in each commit with own hash prefix", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewCritical(
				"fe85083: type must be one of [feat,fix,chore,refactor,perf,test,docs]",
			),
			issue.NewCritical(
				"3dd0e0b: type must be one of [feat,fix,chore,refactor,perf,test,docs]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"fe85083dd422a6ee60443994f1e51d8fbbac0001",
									commit.ParseMessage(
										[]byte("[TASK-123] wrong(ui): add new feature"),
									),
								),
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-546] wrong(ui): add new payment button"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithTypes(
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit without ticket in subject", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewCritical(
				"3dd0e0b: ticket not found in subject by expression '(TASK|PROJ|BUG)-[0-9]+'",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("feat(ui): add new feature"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
								message.NewConfig(
									message.WithSubject(
										subject.NewConfig(
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
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit with not allowed type", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				"3dd0e0b: type must be one of [feat,fix,chore,refactor,perf,test,docs]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-123] wrong(ui): add new feature"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithTypes(
								types.NewConfig(
									value.NewStrings(
										issue.Warning,
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit with empty type escalated to critical", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewCritical(
				"3dd0e0b: type must be one of [feat,fix,chore,refactor,perf,test,docs]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-123] fix payment state machine"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithTypes(
								types.NewConfig(
									value.NewStrings(
										issue.Warning,
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
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit with invalid scope", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				"3dd0e0b: scope not found in subject by expression '^[A-Za-z _-]+$'",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("feat(123): add new feature"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithScope(
								scope.NewConfig(
									value.NewString(
										issue.Warning,
										`^[A-Za-z _-]+$`,
									),
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit with too short unicode description", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewWarning(
				"3dd0e0b: subject description length is not in range [10-72]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("feat(ui): платеж"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
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
										),
									),
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("found commit with too long body", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{
			issue.NewInfo(
				"3dd0e0b: body is not in range [0-5]",
			),
		}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-123] feat(ui): add new feature\n\nbody text"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
								message.NewConfig(
									message.WithBody(
										body.NewConfig(
											body.WithLength(
												value.NewRange(
													issue.Info,
													0,
													5,
												),
											),
										),
									),
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("not found when body present for mandated type", func(t *testing.T) {
		t.Parallel()
		want := []issue.Issue{}
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
					mockGit.WithMergeBase(
						func(_ context.Context, _, _ string) (commit.Hash, error) {
							return commit.NewHash("fe85083dd422a6ee60443994f1e51d8fbbac0001"), nil
						},
					),
					mockGit.WithCommits(
						func(_ context.Context, _ commit.Hash) ([]commit.Commit, error) {
							return []commit.Commit{
								commit.NewCommit(
									"3dd0e0bb904eb0075316f04350b328d06dc825b4",
									commit.ParseMessage(
										[]byte("[TASK-123] fix(ui): add new feature\n\nbody text"),
									),
								),
							}, nil
						},
					),
				),
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
					config.WithCommit(
						commitConfig.NewConfig(
							commitConfig.WithMessage(
								message.NewConfig(
									message.WithBody(
										body.NewConfig(
											body.WithLength(
												value.NewRange(
													issue.Info,
													0,
													255,
												),
											),
											body.WithMandate(
												mandate.NewConfig(
													mandate.WithTypes(
														value.NewStrings(
															issue.Critical,
															"fix",
														),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			).
			Issues(context.Background())
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("failed on empty default branch name", func(t *testing.T) {
		t.Parallel()
		got, gotErr := impl.
			NewLinter(
				mockGit.NewClient(
					mockGit.WithCurrentBranch(
						func(_ context.Context) (branch.Branch, error) {
							return branch.
								NewBranch(
									branch.NewName("feature/TASK-123"),
								), nil
						},
					),
				),
				*config.NewConfig(),
			).
			Issues(context.Background())
		assert.Nil(t, got)
		assert.EqualError(
			t,
			gotErr,
			"invalid config: empty default branch name",
		)
	})
}
