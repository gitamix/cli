# Contributing

Thank you for your interest in this project!
Here’s how you can contribute to our project and help us make it even better.

> [!IMPORTANT]
> Read [STYLE.MD](STYLE.MD) before writing code or opening a Pull Request.
> This repository uses a project-specific code style,
> and contributions are expected to follow it.

## Requirements

Before contributing, follow these rules:

- follow the project code style
- discuss changes through an issue
and submit them through a pull request
- work in a dedicated feature branch
and pass all required checks before submission
- use the required commit message convention
and update the changelog for user-visible changes
- follow the project release process:
  - semantic versioning
  - no manual tags
  - and a version tag for every commit in the master branch

## Getting Started

1. Clone the repository.

    ```bash
    git clone https://github.com/yourusername/your-repo.git
    ```

2. Run the personalization script
to replace template placeholders with your GitHub identity.

    ```bash
    sh ./scripts/itsme.sh <your_github_nickname> <your_github_repo>
    ```

3. [Pick an open issue](https://github.com/Ilya-Sitnikov/gotmpl/issues)
and work on it.

    If there is no suitable issue for your idea,
    open a short issue describing the proposal or contact the maintainer first.

4. Fork the repository and create a feature branch.

    ```bash
    git checkout -b feature/your-feature
    ```

    > [!TIP]
    > If you work from an existing GitHub issue,
    > you can also use the branch name suggested there.

5. Set up the project.

    ```bash
    # Sets Git hooks and registers you in the project templates
    make setup
    ```

6. Commit your changes according to [STYLE.MD](STYLE.MD).

    > [!IMPORTANT]
    > Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/),
    > or your commits may be rejected

7. Check the project before submitting a pull request.

    ```bash
    # Runs tests, linters, and other checks
    make check
    ```

    > [!TIP]
    > Some checks also run on commit, push, and in GitHub Actions.
    > You can still run `make check` locally to catch issues earlier.

8. Update `CHANGELOG.md` for user-visible changes
and keep it in the [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/) format.

    > [!TIP]
    > It is often easier to keep changelog updates in a dedicated commit.

    If you are not sure which version the change belongs to,
    discuss it with the maintainer first.

9. Submit a pull request.

    Use `make pr release <version>` to prepare a release pull request
    from `.github/PULL_REQUEST_TEMPLATE/release.md`.

    ```sh
    make pr release v1.0.0
    ```

    The command reads the matching version block from `CHANGELOG.md`,
    injects it into the PR summary, and opens the GitHub Pull Request page
    in the browser with prefilled query parameters.

    If the requested version is missing from `CHANGELOG.md`, the command fails.

    Your PR must pass all workflows and be approved in code review
    by at least one maintainer before it can be merged.

    > [!IMPORTANT]
    > Merge Pull Requests with squash

    Merging with squash keeps the `master` branch clean
    and avoids polluting it with intermediate
    or unnecessary commits from the review process.

## Versioning and Release

This project uses [Semantic Versioning](https://semver.org/lang/ru/).

Do not create release tags manually.
Use the `Release` GitHub Actions workflow instead:

1. Prepare `CHANGELOG.md` for the future version
in the [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/) format.
2. Run the `Release` workflow
and pass the future version number without the `v` prefix.
3. The workflow validates the version,
extracts the matching section from `CHANGELOG.md`,
creates a `v`-prefixed tag for the latest commit in the `master` branch,
and publishes the GitHub release with those changelog notes.

## Code of Conduct

Be respectful and constructive. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Questions

Open an issue or contact the maintainers directly (see [README](README.md)).
