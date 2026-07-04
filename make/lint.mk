# Lint the codebase.
.PHONY: lint
lint:
	@golangci-lint run ./...
