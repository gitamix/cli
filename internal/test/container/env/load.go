package env

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	testcontainers "github.com/testcontainers/testcontainers-go"
)

// Load loads environment variables from a file inside a Docker container.
//
// It reads the file at the specified filepath inside the container,
// parses its contents, and returns a map of environment variables.
func Load(
	ctx context.Context,
	ctr testcontainers.Container,
	filepath string,
) (Variables, error) {
	rc, err := ctr.CopyFileFromContainer(
		ctx,
		filepath,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rc.Close()
	}()
	bb, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	vars := make(Variables)
	scanner := bufio.NewScanner(
		strings.NewReader(string(bb)),
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("invalid env line: %q", line)
		}
		vars.Set(
			strings.TrimSpace(key),
			strings.TrimSpace(value),
		)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(vars) == 0 {
		return nil, errors.New("no environment variables found in file")
	}
	return vars, nil
}
