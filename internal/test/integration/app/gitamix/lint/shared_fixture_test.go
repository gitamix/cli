//go:build integration
// +build integration

package lint_test

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/gitamix/cli/internal/test/container"
	"github.com/gitamix/cli/internal/test/container/env"
	ctrfx "github.com/gitamix/cli/internal/test/fixture/container"
)

var (
	// fixtureOnce guarantees the shared integration fixture
	// is initialized only once per test process.
	fixtureOnce sync.Once

	// sharedFX stores the package-level container fixture instance
	// reused by integration tests in this package.
	sharedFX *ctrfx.Container

	// sharedErr stores the initialization error returned
	// during the first shared fixture setup attempt.
	sharedErr error
)

// sharedContainerFixture returns the package-level shared container fixture
// used by integration tests in this package.
//
// It initializes the fixture only once and fails the test
// if the shared setup could not be completed successfully.
func sharedContainerFixture(t *testing.T) *ctrfx.Container {
	t.Helper()
	fixtureOnce.Do(func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Minute,
		)
		defer cancel()
		sharedFX, sharedErr = setupSharedContainerFixture(ctx)
	})
	if sharedErr != nil {
		t.Fatalf("failed to setup shared container fixture: %v", sharedErr)
	}
	return sharedFX
}

// setupSharedContainerFixture creates a container fixture backed
// by the shared test container
// and loads the fixture environment variables from it.
func setupSharedContainerFixture(
	ctx context.Context,
) (*ctrfx.Container, error) {
	ctr, err := container.TestContainer(
		ctx,
		2*time.Minute,
		tc.FromDockerfile{
			Context:        mustRepoRoot(),
			Dockerfile:     "docker/git/Dockerfile",
			Repo:           "gitamix/gitamix",
			Tag:            "local",
			KeepImage:      true,
			BuildLogWriter: io.Discard,
		},
	)
	if err != nil {
		return nil, err
	}
	vars, err := env.Load(
		ctx,
		ctr,
		"/opt/fixture/repo/.hashes.fixture.env",
	)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("failed to load env: %w", err)
	}
	return ctrfx.NewContainer(ctr, vars), nil
}

// terminateSharedContainerFixture stops and removes the shared container fixture
// if it has been initialized for the current test process.
func terminateSharedContainerFixture(
	ctx context.Context,
	timeout time.Duration,
) {
	if sharedFX == nil || sharedFX.Container() == nil {
		return
	}
	termCtx, termCancel := context.WithTimeout(ctx, timeout)
	defer termCancel()
	sharedFX.Terminate(termCtx, timeout)
}

// mustRepoRoot returns the root directory of the repository containing this test code
// and panics if it cannot determine the repository root.
//
// This function is useful for locating files relative to the repository root in tests.
func mustRepoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("resolve test file path")
	}
	return filepath.Clean(
		filepath.Join(
			filepath.Dir(file),
			"../../../../../..",
		),
	)
}
