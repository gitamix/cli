#!/usr/bin/env bash
#
# Initializes a deterministic fake Git repository for integration tests.
# This script creates a structured Git history with predictable commits, branches,
# and merge commits that can be used to test Git-related functionality in a
# controlled, reproducible environment.
#
# Usage: ./init_fake_git_tree.sh [repo_directory]
#   repo_directory - Path where the fixture repository will be created (default: /opt/fixture/repo)
#
# The generated repository includes:
# - A main branch with 10 first-parent commits
# - A feature-merge branch with 7 commits that gets merged into main
# - A release/1.0.0 branch with merged alpha and beta release features
# - All commits have deterministic timestamps spanning from 2026-01-01 to 2026-01-23
# - Commit hashes are saved to a .hashes.fixture.env file for test reference
# - A bare clone at /opt/fixture/origin.git serving as the origin remote
# - A final clean branch feature/TASK-123-lint checked out at main's HEAD
#

set -euo pipefail

# Parse arguments: repository directory path (default: /opt/fixture/repo)
repo_dir="${1:-/opt/fixture/repo}"
# File where commit hashes will be stored for test access
hashes_file="${repo_dir}/.hashes.fixture.env"

# Disable system-wide Git configuration to ensure a clean, isolated environment
export GIT_CONFIG_NOSYSTEM=1
# Set HOME to a temporary directory to isolate Git config
export HOME
HOME="$(mktemp -d)"
# Ensure cleanup of temporary directory on script exit
trap 'rm -rf "$HOME"' EXIT

# Create repository directory and remove any existing Git metadata
mkdir -p "$repo_dir"
rm -rf "$repo_dir/.git"

# Initialize new Git repository with 'main' as the default branch
cd "$repo_dir"
git init -b main .
# Configure committer identity for all commits
git config user.name 'Gitamix Fixture'
git config user.email 'fixture@gitamix.local'
# Disable GPG signing for all commits (speeds up testing)
git config commit.gpgSign false

# Helper function to write content to a file, creating parent directories if needed
# Usage: write_file <file_path> <content_lines...>
write_file() {
    local path="$1"
    shift
    mkdir -p "$(dirname "$path")"
    printf '%s\n' "$@" > "$path"
}

# Helper function to create a commit with a specific timestamp
# Usage: commit_with_date <timestamp> <commit_message>
commit_with_date() {
    local when="$1"
    local message="$2"
    GIT_AUTHOR_DATE="$when" \
    GIT_COMMITTER_DATE="$when" \
    git add -A
    GIT_AUTHOR_DATE="$when" \
    GIT_COMMITTER_DATE="$when" \
    git commit --quiet -m "$message"
}

# Helper function to create a merge commit with a specific timestamp
# Usage: merge_with_date <timestamp> <source_branch> <merge_message>
merge_with_date() {
    local when="$1"
    local branch="$2"
    local message="$3"
    GIT_AUTHOR_DATE="$when" \
    GIT_COMMITTER_DATE="$when" \
    git merge --quiet --no-ff --no-edit -m "$message" "$branch"
}

# ============================================================================
# MAIN BRANCH HISTORY (10 first-parent commits)
# ============================================================================

# Commit 1: Initialize repository with basic files
write_file \
    README.md \
    '# Gitamix fixture repository' \
    '' \
    'Deterministic git history for integration tests.'
write_file \
    docs/main/history.md \
    'main history seed'
commit_with_date \
    '2026-01-01T00:00:00+0000' \
    'chore: initialize fixture repository'

# Commit 2: Update history documentation
write_file \
    docs/main/history.md \
    'main history seed' \
    'main-02'
commit_with_date \
    '2026-01-02T00:00:00+0000' \
    'docs(main): add second main commit'

# Commit 3: Add main application file
write_file \
    app/main.txt \
    'main-03'
commit_with_date \
    '2026-01-03T00:00:00+0000' \
    'feat(main): add third main commit'

# Commit 4: Fix/add to main application file
write_file \
    app/main.txt \
    'main-03' \
    'main-04'
commit_with_date \
    '2026-01-04T00:00:00+0000' \
    'fix(main): add fourth main commit'

# Commit 5: Refactor main application file
write_file \
    app/main.txt \
    'main-03' \
    'main-04' \
    'main-05'
commit_with_date \
    '2026-01-05T00:00:00+0000' \
    'refactor(main): add fifth main commit'

# ============================================================================
# FEATURE BRANCH: feature-merge (7 commits, then merged into main)
# ============================================================================

# Create and switch to feature branch
git checkout -b feature-merge

# Feature branch commits 1-7
write_file \
    features/merge.txt \
    'feature-merge-01'
commit_with_date \
    '2026-01-06T00:00:00+0000' \
    'feat(feature-merge): add commit 01'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02'
commit_with_date \
    '2026-01-07T00:00:00+0000' \
    'feat(feature-merge): add commit 02'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02' \
    'feature-merge-03'
commit_with_date \
    '2026-01-08T00:00:00+0000' \
    'feat(feature-merge): add commit 03'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02' \
    'feature-merge-03' \
    'feature-merge-04'
commit_with_date \
    '2026-01-09T00:00:00+0000' \
    'feat(feature-merge): add commit 04'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02' \
    'feature-merge-03' \
    'feature-merge-04' \
    'feature-merge-05'
commit_with_date \
    '2026-01-10T00:00:00+0000' \
    'feat(feature-merge): add commit 05'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02' \
    'feature-merge-03' \
    'feature-merge-04' \
    'feature-merge-05' \
    'feature-merge-06'
commit_with_date \
    '2026-01-11T00:00:00+0000' \
    'feat(feature-merge): add commit 06'

write_file \
    features/merge.txt \
    'feature-merge-01' \
    'feature-merge-02' \
    'feature-merge-03' \
    'feature-merge-04' \
    'feature-merge-05' \
    'feature-merge-06' \
    'feature-merge-07'
commit_with_date \
    '2026-01-12T00:00:00+0000' \
    'feat(feature-merge): add commit 07'

# Merge feature branch back into main
git checkout main
merge_with_date \
    '2026-01-13T00:00:00+0000' \
    'feature-merge' \
    'merge(feature-merge): integrate 7 commits'

# Continue with main branch commits (after merge)
# Commit 6: Update main application file
write_file \
    app/main.txt \
    'main-03' \
    'main-04' \
    'main-05' \
    'main-07'
commit_with_date \
    '2026-01-14T00:00:00+0000' \
    'chore(main): add seventh first-parent commit'

# Commit 7: Add to main application file
write_file \
    app/main.txt \
    'main-03' \
    'main-04' \
    'main-05' \
    'main-07' \
    'main-08'
commit_with_date \
    '2026-01-15T00:00:00+0000' \
    'feat(main): add eighth first-parent commit'

# Commit 8: Fix main application file
write_file \
    app/main.txt \
    'main-03' \
    'main-04' \
    'main-05' \
    'main-07' \
    'main-08' \
    'main-09'
commit_with_date \
    '2026-01-16T00:00:00+0000' \
    'fix(main): add ninth first-parent commit'

# Commit 9: Update documentation in main
write_file \
    app/main.txt \
    'main-03' \
    'main-04' \
    'main-05' \
    'main-07' \
    'main-08' \
    'main-09' \
    'main-10'
commit_with_date \
    '2026-01-17T00:00:00+0000' \
    'docs(main): add tenth first-parent commit'

# ============================================================================
# RELEASE BRANCH: release/1.0.0
# ============================================================================

# Create release branch from main
git checkout -b release/1.0.0

# ============================================================================
# FEATURE BRANCH: feature/release-alpha (2 commits, then merged into release)
# ============================================================================

# Create alpha release feature branch
git checkout -b feature/release-alpha

write_file \
    release/alpha.txt \
    'release-alpha-01'
commit_with_date \
    '2026-01-18T00:00:00+0000' \
    'feat(release-alpha): add commit 01'

write_file \
    release/alpha.txt \
    'release-alpha-01' \
    'release-alpha-02'
commit_with_date \
    '2026-01-19T00:00:00+0000' \
    'feat(release-alpha): add commit 02'

# Merge alpha release into release branch
git checkout release/1.0.0
merge_with_date \
    '2026-01-20T00:00:00+0000' \
    'feature/release-alpha' \
    'merge(feature/release-alpha): prepare release'

# ============================================================================
# FEATURE BRANCH: feature/release-beta (2 commits, then merged into release)
# ============================================================================

# Create beta release feature branch
git checkout -b feature/release-beta

write_file \
    release/beta.txt \
    'release-beta-01'
commit_with_date \
    '2026-01-21T00:00:00+0000' \
    'feat(release-beta): add commit 01'

write_file \
    release/beta.txt \
    'release-beta-01' \
    'release-beta-02'
commit_with_date \
    '2026-01-22T00:00:00+0000' \
    'feat(release-beta): add commit 02'

# Merge beta release into release branch
git checkout release/1.0.0
merge_with_date \
    '2026-01-23T00:00:00+0000' \
    'feature/release-beta' \
    'merge(feature/release-beta): prepare release'

# ============================================================================
# SAVE COMMIT HASHES FOR TEST REFERENCE
# ============================================================================

# Generate a file with commit hashes that tests can use to verify Git operations
# The hashes are saved in a format that can be sourced in shell scripts
cat > "$hashes_file" <<EOF
GITAMIX_FIXTURE_REPO=$repo_dir
MAIN_FIRST_PARENT_01=$(git rev-list --first-parent --reverse main | sed -n '1p')
MAIN_FIRST_PARENT_02=$(git rev-list --first-parent --reverse main | sed -n '2p')
MAIN_FIRST_PARENT_03=$(git rev-list --first-parent --reverse main | sed -n '3p')
MAIN_FIRST_PARENT_04=$(git rev-list --first-parent --reverse main | sed -n '4p')
MAIN_FIRST_PARENT_05=$(git rev-list --first-parent --reverse main | sed -n '5p')
MAIN_FIRST_PARENT_06=$(git rev-list --first-parent --reverse main | sed -n '6p')
MAIN_FIRST_PARENT_07=$(git rev-list --first-parent --reverse main | sed -n '7p')
MAIN_FIRST_PARENT_08=$(git rev-list --first-parent --reverse main | sed -n '8p')
MAIN_FIRST_PARENT_09=$(git rev-list --first-parent --reverse main | sed -n '9p')
MAIN_FIRST_PARENT_10=$(git rev-list --first-parent --reverse main | sed -n '10p')
FEATURE_MERGE_01=$(git rev-list --reverse feature-merge ^main~5 | sed -n '1p')
FEATURE_MERGE_02=$(git rev-list --reverse feature-merge ^main~5 | sed -n '2p')
FEATURE_MERGE_03=$(git rev-list --reverse feature-merge ^main~5 | sed -n '3p')
FEATURE_MERGE_04=$(git rev-list --reverse feature-merge ^main~5 | sed -n '4p')
FEATURE_MERGE_05=$(git rev-list --reverse feature-merge ^main~5 | sed -n '5p')
FEATURE_MERGE_06=$(git rev-list --reverse feature-merge ^main~5 | sed -n '6p')
FEATURE_MERGE_07=$(git rev-list --reverse feature-merge ^main~5 | sed -n '7p')
RELEASE_HEAD=$(git rev-parse release/1.0.0)
RELEASE_ALPHA_HEAD=$(git rev-parse feature/release-alpha)
RELEASE_BETA_HEAD=$(git rev-parse feature/release-beta)
EOF

# ============================================================================
# CLEAN LINT BRANCH: feature/TASK-123-lint
# ============================================================================

# Create a clean branch checked out at main's HEAD for lint integration tests.
# The branch name matches the pattern ^(feature|bugfix|hotfix)/[A-Z]+-\d+
# and contains the task identifier TASK-123.
# It points at the same commit as main, so it introduces no commits
# that are not pushed to its remote counterpart yet.
git checkout --quiet -b feature/TASK-123-lint

# ============================================================================
# REMOTE ORIGIN
# ============================================================================

# Create a bare clone of the fixture repository acting as the origin remote.
git clone --quiet --bare "$repo_dir" /opt/fixture/origin.git
# Register the bare clone as the origin remote of the fixture repository.
git remote add origin /opt/fixture/origin.git
# Fetch from origin to create remote-tracking refs (refs/remotes/origin/*).
git fetch --quiet origin
