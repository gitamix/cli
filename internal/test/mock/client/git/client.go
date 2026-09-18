package git

import (
	"context"

	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
)

// MockGit simulates the behavior of the git command-line client in tests.
//
// MockGit implements the command-line interface abstracted by the client
// of the internal/client/git package. Every method delegates the call
// to a function specified at creation and panics if no function
// is specified for it.
type MockGit struct {
	// funcs holds the behaviors of the mock methods.
	//
	// Each function is executed when the corresponding
	// method of MockGit is called.
	funcs funcs
}

// NewClient creates a new MockGit instance with the provided options.
//
// Each option assigns the behavior of the corresponding mock method.
// Methods without an assigned behavior panic on call.
func NewClient(opts ...Option) *MockGit {
	m := &MockGit{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// CurrentBranch simulates the behavior of the CurrentBranch method.
//
// It returns the branch and the error produced by the function
// specified with the WithCurrentBranch option.
//
// Panics if the behavior for the CurrentBranch method is not specified.
func (m *MockGit) CurrentBranch(
	ctx context.Context,
) (branch.Branch, error) {
	if m.funcs.CurrentBranch == nil {
		panic("not specified behavior for CurrentBranch method")
	}
	return m.funcs.CurrentBranch(ctx)
}

// MergeBase simulates the behavior of the MergeBase method.
//
// It returns the commit hash and the error produced by the function
// specified with the WithMergeBase option.
//
// Panics if the behavior for the MergeBase method is not specified.
func (m *MockGit) MergeBase(
	ctx context.Context,
	target string,
	curr string,
) (commit.Hash, error) {
	if m.funcs.MergeBase == nil {
		panic("not specified behavior for MergeBase method")
	}
	return m.funcs.MergeBase(ctx, target, curr)
}
