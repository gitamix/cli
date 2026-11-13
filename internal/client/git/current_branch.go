package git

import (
	"context"

	"github.com/gitamix/types/branch"
)

// CurrentBranch returns the current branch of the repository
// or an error if the git command execution fails.
func (c *Client) CurrentBranch(
	ctx context.Context,
) (branch.Branch, error) {
	return c.cli.CurrentBranch(ctx)
}
