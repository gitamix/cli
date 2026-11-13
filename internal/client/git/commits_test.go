package git_test

import (
	"context"
	"testing"

	"github.com/gitamix/types/commit"
	"github.com/stretchr/testify/assert"

	impl "github.com/gitamix/cli/internal/client/git"
	gitmk "github.com/gitamix/cli/internal/test/mock/client/git"
)

func TestClient_Commits(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		want := []commit.Commit{
			commit.NewCommit(
				commit.NewHash("89abcdef"),
				commit.NewMessage(
					commit.NewSubject(
						commit.NewType("fix"),
						commit.NewScope("ui"),
						commit.NewDescription("resolve button alignment issue"),
					),
					commit.NewBody(
						[]byte("Fixed the alignment issue of the button in the UI."),
					),
				),
			),
			commit.NewCommit(
				commit.NewHash("1234567"),
				commit.NewMessage(
					commit.NewSubject(
						commit.NewType("feat"),
						commit.NewScope("domain"),
						commit.NewDescription("add new feature"),
					),
					commit.NewBody(
						[]byte("Added a new feature to the domain layer."),
					),
				),
			),
		}
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCommits(
					func(
						_ context.Context,
						_ commit.Hash,
					) ([]commit.Commit, error) {
						return []commit.Commit{
							commit.NewCommit(
								commit.NewHash("89abcdef"),
								commit.NewMessage(
									commit.NewSubject(
										commit.NewType("fix"),
										commit.NewScope("ui"),
										commit.NewDescription("resolve button alignment issue"),
									),
									commit.NewBody(
										[]byte("Fixed the alignment issue of the button in the UI."),
									),
								),
							),
							commit.NewCommit(
								commit.NewHash("1234567"),
								commit.NewMessage(
									commit.NewSubject(
										commit.NewType("feat"),
										commit.NewScope("domain"),
										commit.NewDescription("add new feature"),
									),
									commit.NewBody(
										[]byte("Added a new feature to the domain layer."),
									),
								),
							),
						}, nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.Commits(
			ctx,
			commit.NewHash("1234567"),
		)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})

	t.Run("failed to get commits", func(t *testing.T) {
		t.Parallel()
		var want []commit.Commit
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCommits(
					func(
						_ context.Context,
						_ commit.Hash,
					) ([]commit.Commit, error) {
						return nil, assert.AnError
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.Commits(
			ctx,
			commit.NewHash("1234567"),
		)
		assert.Equal(t, want, got)
		assert.ErrorIs(t, gotErr, assert.AnError)
	})

	t.Run("got empty commits", func(t *testing.T) {
		t.Parallel()
		want := []commit.Commit{}
		git := impl.NewClient(
			gitmk.NewClient(
				gitmk.WithCommits(
					func(
						_ context.Context,
						_ commit.Hash,
					) ([]commit.Commit, error) {
						return []commit.Commit{}, nil
					},
				),
			),
		)
		ctx := context.Background()
		got, gotErr := git.Commits(
			ctx,
			commit.NewHash("1234567"),
		)
		assert.Equal(t, want, got)
		assert.NoError(t, gotErr)
	})
}
