package container

import (
	"context"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestContainer inits and returns test container with Git fixture image
// built from Dockerfile in the repository root.
//
// It uses the provided context and timeout for container initialization.
func TestContainer(
	ctx context.Context,
	timeout time.Duration,
	dockf tc.FromDockerfile,
) (tc.Container, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		timeout,
	)
	defer cancel()
	ctr, err := tc.Run(
		ctx,
		"",
		tc.WithDockerfile(dockf),
		tc.WithWaitStrategy(
			wait.
				ForExec([]string{"git", "rev-parse", "main"}).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}
	return ctr, nil
}
