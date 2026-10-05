package lint

import (
	"context"
	"fmt"

	lintBranch "github.com/gitamix/lint/branch"
	lintCommit "github.com/gitamix/lint/commit"
	"github.com/gitamix/lint/issue"

	"github.com/gitamix/cli/internal/validation/gitamix/lint/config"
)

// Issues returns all issues found in the current branch name
// and in the commit messages introduced by that branch.
//
// Commit messages are checked for every commit reachable
// from the current HEAD but not from the merge base
// of the current branch with the default branch of the configuration,
// that is, for the commits that are not merged to the default branch yet.
//
// An empty default branch name in the configuration aborts the linting
// with an error, since the merge base cannot be resolved then.
//
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
	err = config.
		NewValidation(l.cfg).
		Validate()
	if err != nil {
		return nil, fmt.Errorf(
			"invalid config: %w",
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
		l.cfg.Branch().Default().Name(),
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
