package git_test

import (
	"context"
	"testing"

	"github.com/gitamix/types/commit"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/client/git"
	gitmk "github.com/gitamix/cli/internal/test/mock/client/git"
)

func TestClient_MergeBase(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		want := commit.NewHash(
			"1a2b3c4d5e6f7a8b",
		)
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithMergeBase(
					func(
						_ context.Context,
						_ string,
						_ string,
					) (commit.Hash, error) {
						return commit.NewHash(
							"1a2b3c4d5e6f7a8b",
						), nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.MergeBase(
			ctx,
			"master",
			"feature",
		)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("failed to get merge base", func(t *testing.T) {
		t.Parallel()
		var want commit.Hash
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithMergeBase(
					func(
						_ context.Context,
						_ string,
						_ string,
					) (commit.Hash, error) {
						return "", assert.AnError
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.MergeBase(
			ctx,
			"master",
			"feature",
		)
		assert.Equal(t, want, got)
		assert.ErrorIs(t, gotErr, assert.AnError)
	})

	t.Run("got empty merge base hash", func(t *testing.T) {
		t.Parallel()
		want := commit.NewHash(
			"",
		)
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithMergeBase(
					func(
						_ context.Context,
						_ string,
						_ string,
					) (commit.Hash, error) {
						return commit.NewHash(
							"",
						), nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.MergeBase(
			ctx,
			"master",
			"feature",
		)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})
}
