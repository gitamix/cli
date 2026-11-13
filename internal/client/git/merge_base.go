package git

import (
	"context"

	"github.com/gitamix/types/commit"
)

// MergeBase returns the merge base commit hash of the two specified
// commits or branches or an error if the git command execution fails.
func (c *Client) MergeBase(
	ctx context.Context,
	target string,
	curr string,
) (commit.Hash, error) {
	return c.cli.MergeBase(ctx, target, curr)
}
