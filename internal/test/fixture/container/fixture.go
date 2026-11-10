package container

import (
	"context"
	"time"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/gitamix/cli/internal/test/container/env"
)

// Container represents a test container fixture
// that is used for testing purposes.
type Container struct {
	// ctr is the test Docker container
	// we use to run the tests against.
	ctr tc.Container

	// env is the collection of environment variables
	// that contain commit hashes.
	env env.Variables
}

// NewContainer creates a new Container instance
// with the given Docker container and commit hashes.
func NewContainer(
	ctr tc.Container,
	env env.Variables,
) *Container {
	return &Container{
		ctr: ctr,
		env: env,
	}
}

// Container returns the test Docker container
// associated with the actual fixture.
func (fx *Container) Container() tc.Container {
	return fx.ctr
}

// Env returns the collection of environment variables
// that contain commit hashes.
func (fx *Container) Env() env.Variables {
	return fx.env
}

// Terminate stops and removes the test Docker container
// associated with the actual fixture, using the specified timeout.
func (fx *Container) Terminate(
	ctx context.Context,
	timeout time.Duration,
) error {
	if fx.ctr == nil {
		return nil
	}
	return fx.ctr.Terminate(
		ctx,
		tc.StopTimeout(timeout),
	)
}
