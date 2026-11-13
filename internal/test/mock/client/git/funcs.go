package git

import (
	"context"

	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
)

// funcs holds the functions that are executed when the corresponding
// methods of MockGit are called.
type funcs struct {
	// CurrentBranch is the function that is executed
	// when CurrentBranch of MockGit is called.
	CurrentBranch func(
		ctx context.Context,
	) (branch.Branch, error)
	// MergeBase is the function that is executed
	// when MergeBase of MockGit is called.
	MergeBase func(
		ctx context.Context,
		target string,
		curr string,
	) (commit.Hash, error)
	// Commits is the function that is executed
	// when Commits of MockGit is called.
	Commits func(
		ctx context.Context,
		hash commit.Hash,
	) ([]commit.Commit, error)
}
