package git

import (
	"context"

	"github.com/gitamix/types/branch"
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
