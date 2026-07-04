#!/bin/sh

set -eu

base_branch="master"

fail() {
	echo "$1" >&2
	exit 1
}

usage() {
	echo "Usage: make pr release <version>" >&2
}

require_command() {
	command_name="$1"
	message="$2"

	if ! command -v "$command_name" >/dev/null 2>&1; then
		fail "$message"
	fi
}

resolve_repo_url() {
	remote_url=$(git remote get-url origin 2>/dev/null || true)

	if [ -z "$remote_url" ]; then
		fail "Unable to detect git remote 'origin'"
	fi

	case "$remote_url" in
		https://github.com/*)
			repo_url=${remote_url%.git}
			;;
		git@github.com:*)
			repo_path=${remote_url#git@github.com:}
			repo_url="https://github.com/${repo_path%.git}"
			;;
		ssh://git@github.com/*)
			repo_path=${remote_url#ssh://git@github.com/}
			repo_url="https://github.com/${repo_path%.git}"
			;;
		*)
			fail "Unsupported git remote URL for origin: $remote_url"
			;;
	esac
	
}

escape_sed_replacement() {
	printf '%s' "$1" | sed 's/[\\/&]/\\&/g'
}

url_encode() {
	perl -0pe 's/([^A-Za-z0-9_.~-])/sprintf("%%%02X", ord($1))/ge'
}

parse_args() {
	if [ "$#" -ne 2 ] || [ "$1" != "release" ] || [ -z "${2:-}" ]; then
		usage
		exit 1
	fi

	pr_type="release"
	version="$2"
	normalized_version=${version#v}
	version_label="v$normalized_version"
}

prepare_environment() {
	repo_root=$(git rev-parse --show-toplevel)
	template_dir="$repo_root/.github/PULL_REQUEST_TEMPLATE"
	template_path="$template_dir/$pr_type.md"
	changelog_path="$repo_root/CHANGELOG.md"
	current_branch=$(git branch --show-current)

	if [ ! -f "$template_path" ]; then
		fail "PR template not found: $template_path"
	fi

	if [ ! -f "$changelog_path" ]; then
		fail "CHANGELOG.md not found: $changelog_path"
	fi

	if [ -z "$current_branch" ]; then
		fail "Unable to detect current git branch"
	fi

	require_command perl "perl is required to build the GitHub PR URL"
	resolve_repo_url
}

extract_release_summary() {
	summary_file=$(mktemp)
	trap 'rm -f "$temp_file" "$summary_file"' EXIT HUP INT TERM

	awk -v version="$normalized_version" '
		BEGIN {
			in_block = 0
			found = 0
		}
		$0 ~ "^## \\[" version "\\]" {
			in_block = 1
			found = 1
			started = 0
			next
		}
		in_block && $0 ~ "^## " {
			exit
		}
		in_block {
			if (!started && $0 ~ /^[[:space:]]*$/) {
				next
			}
			started = 1
			print
		}
		END {
			if (!found) {
				exit 2
			}
		}
	' "$changelog_path" > "$summary_file" || exit_code=$?

	if [ "${exit_code:-0}" -eq 2 ]; then
		fail "Release notes for version $version_label were not found in CHANGELOG.md"
	fi

	if [ "${exit_code:-0}" -ne 0 ]; then
		fail "Failed to read release notes from CHANGELOG.md"
	fi

	if ! grep -q '[^[:space:]]' "$summary_file"; then
		fail "Release notes for version $version_label are empty in CHANGELOG.md"
	fi
}

render_template() {
	pr_title="release: $version_label"
	awk -v version="$version_label" -v summary_file="$summary_file" '
		function print_summary(path, line) {
			while ((getline line < path) > 0) {
				print line
			}
			close(path)
		}
		{
			gsub(/\{\{VERSION\}\}/, version)
			if ($0 == "{{SUMMARY}}") {
				print_summary(summary_file)
				next
			}
			print
		}
	' "$template_path" > "$temp_file"
}

open_pull_request_page() {
	encoded_base=$(printf '%s' "$base_branch" | url_encode)
	encoded_head=$(printf '%s' "$current_branch" | url_encode)
	encoded_title=$(printf '%s' "$pr_title" | url_encode)
	encoded_body=$(url_encode < "$temp_file")
	compare_url="$repo_url/compare/$encoded_base...$encoded_head?quick_pull=1&title=$encoded_title&body=$encoded_body"

	echo "Opening PR page: $compare_url"
	"$repo_root/scripts/pr/open.sh" "$compare_url"
}

temp_file=$(mktemp)
trap 'rm -f "$temp_file"' EXIT HUP INT TERM

parse_args "$@"
prepare_environment
extract_release_summary
render_template

open_pull_request_page