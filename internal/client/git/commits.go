package git

import (
	"context"

	"github.com/gitamix/types/commit"
)

// Commits returns the commits reachable from the current HEAD
// that are not reachable from the specified hash
// or an error if the git command execution fails.
func (c *Client) Commits(
	ctx context.Context,
	hash commit.Hash,
) ([]commit.Commit, error) {
	return c.cli.Commits(ctx, hash)
}
