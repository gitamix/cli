# Run all tests in the project.
.PHONY: test
test: unit-test

# Run unit tests.
.PHONY: unit-test
unit-test:
	@echo 🧪 Running unit tests...
	@mkdir -p tmp
	@pkgs=$$(go list ./... | grep -v 'internal/test/integration' | grep -v 'internal/test/fake'); \
	if [ -n "$$pkgs" ]; then \
		go test -race -count=1 \
		-covermode=atomic \
		-coverprofile=../tmp/coverage_unit.out \
		$$pkgs; \
	else \
		echo "no packages to test"; \
	fi \
