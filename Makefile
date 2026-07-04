.DEFAULT_GOAL = check

-include .github/ci.env
include make/setup.mk
include make/test.mk
include make/coverage.mk
include make/lint.mk
include make/pr.mk

# Run all checks required to validate the codebase before merging.
.PHONY: check
check: test coverage-check lint 
