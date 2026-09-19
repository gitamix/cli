package lint

import (
	"context"

	"github.com/gitamix/lint/config"
	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
)

// gitClient abstracts the git repository access needed for linting.
//
// The interface is satisfied by the git client adapter
// and keeps the Linter independent of the concrete git implementation.
type gitClient interface {
	// CurrentBranch returns the current branch of the repository
	// or an error if the git operation fails.
	CurrentBranch(
		ctx context.Context,
	) (branch.Branch, error)
	// MergeBase returns the merge base commit hash
	// of the two specified commits or branches
	// or an error if the git operation fails.
	MergeBase(
		ctx context.Context,
		target string,
		curr string,
	) (commit.Hash, error)
	// Commits returns the commits reachable from the current HEAD
	// that are not reachable from the specified hash
	// or an error if the git operation fails.
	Commits(
		ctx context.Context,
		hash commit.Hash,
	) ([]commit.Commit, error)
}

// Linter aggregates issues found in the git branch name
// and commit messages of the user's project.
//
// Linter reads the current repository state through the git client
// and validates it against the rules of the lint configuration.
type Linter struct {
	// git provides access to the git repository of the user's project.
	git gitClient
	// cfg stores the lint rules for branch names and commit messages.
	cfg config.Config
}

// NewLinter creates a new Linter instance
// with the provided git client and lint configuration.
func NewLinter(git gitClient, cfg config.Config) Linter {
	return Linter{
		git: git,
		cfg: cfg,
	}
}
