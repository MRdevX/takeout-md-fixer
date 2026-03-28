#!/usr/bin/env bash
# Shared helpers for release versioning and changelogs (sourced by other scripts).
set -euo pipefail

# Latest tag matching v*.*.* or empty.
latest_version_tag() {
	git fetch --tags --quiet 2>/dev/null || true
	git tag -l 'v*.*.*' --sort=-v:refname | head -1
}

# Strip leading v from v1.2.3 -> 1.2.3
strip_v() {
	local t="$1"
	echo "${t#v}"
}

# Bump patch: 1.2.3 -> 1.2.4
bump_patch() {
	local ver="$1"
	local major minor patch
	IFS='.' read -r major minor patch <<<"${ver//[^0-9.]/}"
	patch=$((patch + 1))
	echo "${major}.${minor}.${patch}"
}

# Next release version string without v prefix (e.g. 0.0.2)
next_release_version() {
	local last tag_ver
	last="$(latest_version_tag)"
	if [[ -z "${last}" ]]; then
		echo "0.0.1"
		return
	fi
	tag_ver="$(strip_v "${last}")"
	bump_patch "${tag_ver}"
}

# Changelog markdown: commits since last v* tag (exclusive) to HEAD (inclusive).
changelog_since_last_release() {
	local prev
	prev="$(latest_version_tag)"
	if [[ -z "${prev}" ]]; then
		git log -100 --pretty=format:'- %s (`%h`)'
	else
		git log "${prev}..HEAD" --pretty=format:'- %s (`%h`)'
	fi
}

# Apply version to build/config.yml (info.version line).
apply_version_to_config() {
	local root="${1:?root}"
	local ver="${2:?version}"
	local f="${root}/build/config.yml"
	if [[ "$(uname -s)" == Darwin ]]; then
		sed -i '' "s/^  version: \".*\"/  version: \"${ver}\"/" "${f}"
	else
		sed -i "s/^  version: \".*\"/  version: \"${ver}\"/" "${f}"
	fi
}
