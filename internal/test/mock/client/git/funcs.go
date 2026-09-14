package git

import (
	"context"

	"github.com/gitamix/types/branch"
)

// funcs holds the functions that are executed when the corresponding
// methods of MockGit are called.
type funcs struct {
	// CurrentBranch is the function that is executed
	// when CurrentBranch of MockGit is called.
	CurrentBranch func(
		ctx context.Context,
	) (branch.Branch, error)
}
