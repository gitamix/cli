# gitamix

[![Go Reference](https://pkg.go.dev/badge/github.com/gitamix/cli.svg)](https://pkg.go.dev/github.com/gitamix/cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gitamix/cli)
[![Release](https://img.shields.io/github/v/release/gitamix/cli?style=flat)](https://github.com/gitamix/cli/releases)
[![Tests](https://github.com/gitamix/cli/actions/workflows/tests.yml/badge.svg)](https://github.com/gitamix/cli/actions/workflows/tests.yml)
[![Quality](https://github.com/gitamix/cli/actions/workflows/quality.yml/badge.svg)](https://github.com/gitamix/cli/actions/workflows/quality.yml)

A CLI tool that improves your git workflow with
[gitamix](https://github.com/gitamix) tools,
such as git commit message and branch name linting.

> [!NOTE]
> The project is a CLI tool built on top of the gitamix libraries:
> [gitamix/git](https://github.com/gitamix/git) fetches the current branch
> and its commits through a pluggable shell,
> [gitamix/lint](https://github.com/gitamix/lint) validates them
> against rules from a YAML config,
> and [gitamix/types](https://github.com/gitamix/types)
> provides the shared domain types.

## Features

- **One-command linting** — `gitamix lint` lints the current branch
  and every commit after the merge base
  without parsing git output yourself.
- **YAML-driven configuration** — declare rules in a single
  `gitamix.lint.yml`; every key is optional,
  and rules are strictly opt-in:
  an empty config reports no issues.
- **Per-rule severity** — every rule is paired with an issue level:
  `critical`, `warning`, or `info`.
  Only `critical` issues fail the run.
- **CI and hooks friendly** — the command exits with a non-zero code
  when critical issues are found,
  so you can wire it into pre-commit hooks and pipelines.

## Installation

Run the command in terminal:

```sh
go install github.com/gitamix/cli@latest
```

All runtime dependencies are pulled in automatically.

## Getting started

Describe your rules in a `gitamix.lint.yml` file.
Every key is optional — only the rules you configure are enforced.

The full configuration reference,
including commit scope, subject task ID, and body mandate rules,
is available in the
[gitamix/lint](https://github.com/gitamix/lint#getting-started) README.

```yaml
# Git branch rules.
branch:
  name:
    issue:
      # Level for branch name rules: critical, warning, info
      level: "warning"
    # Regular expression pattern for branch names.
    pattern: ^(master|release/v\d+\.\d+\.\d+|feature/.+)$
  default:
    name: "master"
# Git commit rules.
commit:
  message:
    subject:
      description:
        length:
          issue:
            # Level for description length rules: critical, warning, info
            level: warning
          # Minimum number of characters allowed
          # in the subject description.
          min: 3
          # Maximum number of characters allowed
          # in the subject description.
          max: 72
    body:
      length:
        issue:
          # Level for body length rules: critical, warning, info
          level: info
          # Minimum number of characters allowed in the body.
          min: 10
          # Maximum number of characters allowed in the body.
          max: 1024
  # Allowed commit types (e.g. feat, fix, chore).
  types:
    issue:
      # Level for commit type rules: critical, warning, info
      level: critical
    list:
      - feat
      - fix
      - chore
      - refactor
      - perf
      - test
      - docs
```

Lint the current branch and its commits:

```sh
gitamix lint --config=/path/to/gitamix.lint.yml
```

Output:

```text
3 issues
critical: abc1234: commit type 'feat!' is not allowed by the permitted types [feat fix chore refactor perf test docs]
warning: abc1234: subject description length is not in range [3-72]
info: abc1234: body length is not in range [10-1024]
Found 1 issues that fail the linting.
```

## Documentation

Full API reference is available on [pkg.go.dev](https://pkg.go.dev/github.com/gitamix/cli)

## Requirements

- [Go](https://go.dev/) 1.25.7 or later.
- Git — the `git` binary must be resolvable in the shell,
  as commands are executed on your behalf.
- A lint configuration file passed via `--config`
  (see the example above).

## Contributing

Want to contribute?
Read [CONTRIBUTING.md](CONTRIBUTING.md)
for the full workflow, repository requirements, and Pull Request process.

Please open an issue to discuss large or breaking changes before implementing.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Author / Contact

Maintained by [Ilya Sitnikov](https://github.com/gitamix)
