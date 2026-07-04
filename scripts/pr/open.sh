#!/bin/sh

set -eu

url="${1:-}"

if [ -z "$url" ]; then
	echo "Usage: scripts/pr/open.sh <url>" >&2
	exit 1
fi

os_name=$(uname -s)

case "$os_name" in
	Darwin)
		launcher="open"
		if ! command -v "$launcher" >/dev/null 2>&1; then
			echo "open is required to launch the GitHub PR page on macOS" >&2
			exit 1
		fi
		exec "$launcher" "$url"
		;;
	Linux)
		launcher="xdg-open"
		if ! command -v "$launcher" >/dev/null 2>&1; then
			echo "xdg-open is required to launch the GitHub PR page on Linux" >&2
			exit 1
		fi
		exec "$launcher" "$url"
		;;
	MINGW*|MSYS*|CYGWIN*)
		if ! command -v cmd.exe >/dev/null 2>&1; then
			echo "cmd.exe is required to launch the GitHub PR page on Windows" >&2
			exit 1
		fi
		exec cmd.exe /c start "" "$url"
		;;
	*)
		echo "Unsupported OS: $os_name" >&2
		exit 1
		;;
esac