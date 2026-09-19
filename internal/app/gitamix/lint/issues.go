package lint

import (
	"context"
	"fmt"

	lintBranch "github.com/gitamix/lint/branch"
	lintCommit "github.com/gitamix/lint/commit"
	"github.com/gitamix/lint/issue"
)

// Issues returns all issues found in the current branch name
// and in the commit messages introduced by that branch.
//
// Commit messages are checked for every commit reachable
// from the current HEAD but not from the merge base
// of the current branch with its remote counterpart on origin,
// that is, for the commits that are not pushed to origin yet.
// A failed git operation aborts the linting
// and returns an error describing the failure.
func (l Linter) Issues(ctx context.Context) ([]issue.Issue, error) {
	currbr, err := l.git.CurrentBranch(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get current branch: %s",
			err,
		)
	}
	iss := make([]issue.Issue, 0, 10)
	iss = append(
		iss,
		lintBranch.
			NewBranch(
				currbr,
				l.cfg.Branch(),
			).
			Issues()...,
	)
	base, err := l.git.MergeBase(
		ctx,
		"origin/"+currbr.String(),
		currbr.String(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get merge base: %s",
			err,
		)
	}
	commits, err := l.git.Commits(
		ctx,
		base,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get commits: %s",
			err,
		)
	}
	iss = append(
		iss,
		lintCommit.
			NewCommits(
				commits,
				l.cfg.Commit(),
			).
			Issues()...,
	)
	return iss, nil
}
