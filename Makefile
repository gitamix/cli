.DEFAULT_GOAL = check

-include .github/ci.env
include make/setup.mk
include make/test.mk
include make/coverage.mk
include make/lint.mk
include make/pr.mk

# Run all checks required to validate the codebase before merging.
.PHONY: check
check: test coverage lint 

# Build a portable binary into bin/.
.PHONY: build
build:
	@echo 🔨 Building binary...
	@mkdir -p bin
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/gitamix .

	@echo ✅ Built bin/gitamix