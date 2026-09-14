package git

import (
	"context"

	"github.com/gitamix/types/branch"
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
