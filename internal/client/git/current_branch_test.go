package git_test

import (
	"context"
	"testing"

	"github.com/gitamix/types/branch"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/client/git"
	gitmk "github.com/gitamix/cli/internal/test/mock/client/git"
)

func TestClient_CurrentBranch(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		want := branch.NewBranch(
			branch.NewName("main"),
		)
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCurrentBranch(
					func(_ context.Context) (branch.Branch, error) {
						return branch.NewBranch(
							branch.NewName("main"),
						), nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.CurrentBranch(ctx)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("failed to get current branch", func(t *testing.T) {
		t.Parallel()
		var want branch.Branch
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCurrentBranch(
					func(_ context.Context) (branch.Branch, error) {
						return branch.Branch{}, assert.AnError
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.CurrentBranch(ctx)
		assert.Equal(t, want, got)
		assert.ErrorIs(t, gotErr, assert.AnError)
	})

	t.Run("got branch with empty name", func(t *testing.T) {
		t.Parallel()
		want := branch.NewBranch(
			branch.NewName(""),
		)
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCurrentBranch(
					func(_ context.Context) (branch.Branch, error) {
						return branch.NewBranch(
							branch.NewName(""),
						), nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.CurrentBranch(ctx)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})
}
