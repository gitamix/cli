package git

import (
	"context"

	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
)

// Option defines a functional option to configure a MockGit instance
// on its creation.
type Option func(*MockGit)

// WithCurrentBranch sets the function that is executed
// when CurrentBranch of MockGit is called.
//
// The specified function determines the branch and the error
// returned by the CurrentBranch method of the mock.
func WithCurrentBranch(
	fn func(
		ctx context.Context,
	) (branch.Branch, error),
) Option {
	return func(r *MockGit) {
		r.funcs.CurrentBranch = fn
	}
}

// WithMergeBase sets the function that is executed
// when MergeBase of MockGit is called.
//
// The specified function determines the commit hash and the error
// returned by the MergeBase method of the mock.
func WithMergeBase(
	fn func(
		ctx context.Context,
		target string,
		curr string,
	) (commit.Hash, error),
) Option {
	return func(r *MockGit) {
		r.funcs.MergeBase = fn
	}
}
