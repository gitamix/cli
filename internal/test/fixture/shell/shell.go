package shell

import (
	"context"
	"fmt"
	"io"

	"github.com/sitnikovik/osxec/command"
	proc "github.com/sitnikovik/osxec/process/execution"
	tc "github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

// Shell is a struct that implements the process.Shell interface
// and executes commands inside a Docker container.
type Shell struct {
	// ctr is the test Docker container in which commands will be executed.
	ctr tc.Container

	// dir is the working directory inside the container
	// where commands will be executed.
	dir string
}

// NewShell creates a new Shell instance that executes commands
// inside the specified Docker container and working directory.
func NewShell(
	ctr tc.Container,
	wkdir string,
) *Shell {
	return &Shell{
		ctr: ctr,
		dir: wkdir,
	}
}

// Execution executes the given command and returns its execution result.
//
// It runs the command inside the specified Docker container and working directory.
func (s *Shell) Execution(
	ctx context.Context,
	cmd command.Command,
) proc.Execution {
	argv := append([]string{cmd.Name()}, cmd.Args()...)
	exitCode, reader, err := s.ctr.Exec(
		ctx,
		argv,
		tcexec.WithWorkingDir(s.dir),
		tcexec.Multiplexed(),
	)
	if err != nil {
		return proc.NewExecution(nil, err)
	}
	bb, readErr := io.ReadAll(reader)
	if readErr != nil {
		return proc.NewExecution(nil, readErr)
	}
	if exitCode != 0 {
		return proc.NewExecution(
			bb,
			fmt.Errorf("exit status %d", exitCode),
		)
	}
	return proc.NewExecution(bb, nil)
}
