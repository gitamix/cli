package git

import (
	"context"

	"github.com/gitamix/types/branch"
	"github.com/gitamix/types/commit"
)

// cli abstracts the command-line interface of git.
//
// Implementations of the interface execute the actual git commands
// and translate their results into domain types.
type cli interface {
	// CurrentBranch returns the current branch of the repository
	// or an error if the git command execution fails.
	CurrentBranch(
		ctx context.Context,
	) (branch.Branch, error)
	// MergeBase returns the merge base commit hash of the two specified
	// commits or branches or an error if the git command execution fails.
	MergeBase(
		ctx context.Context,
		target string,
		curr string,
	) (commit.Hash, error)
	// Commits returns the commits reachable from the current HEAD
	// that are not reachable from the specified hash
	// or an error if the git command execution fails.
	Commits(
		ctx context.Context,
		hash commit.Hash,
	) ([]commit.Commit, error)
}

// Client interacts with the git repository of the user's project.
//
// Client is a thin wrapper: every method delegates the work
// to the underlying command-line interface.
type Client struct {
	// cli executes the git commands behind the Client methods.
	cli cli
}

// NewClient creates a new Client instance
// with the provided command-line interface.
func NewClient(cli cli) *Client {
	return &Client{
		cli: cli,
	}
}
