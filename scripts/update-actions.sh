#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WORKFLOWS_DIR="$REPO_ROOT/{{ .ProjectKebab }}/.github/workflows"

if ! command -v gh &>/dev/null; then
	echo "error: gh CLI is required"
	exit 1
fi

# Extract unique action references (org/repo@vN) from all workflow files
declare -A actions
while IFS= read -r line; do
	# Match "uses: org/repo@vN" patterns, skip local references (./)
	if [[ "$line" =~ uses:\ +([a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+)@v([0-9]+) ]]; then
		repo="${BASH_REMATCH[1]}"
		current="v${BASH_REMATCH[2]}"
		actions["$repo"]="$current"
	fi
done < <(cat "$WORKFLOWS_DIR"/*.yml)

echo "Checking ${#actions[@]} actions for updates..."
echo ""

updated=0
for repo in "${!actions[@]}"; do
	current="${actions[$repo]}"

	# Get the latest release tag from GitHub
	latest_tag=$(gh api "repos/$repo/releases/latest" --jq '.tag_name' 2>/dev/null || echo "")

	if [ -z "$latest_tag" ]; then
		# Some actions don't use releases, try tags instead
		latest_tag=$(gh api "repos/$repo/tags?per_page=1" --jq '.[0].name' 2>/dev/null || echo "")
	fi

	if [ -z "$latest_tag" ]; then
		echo "  SKIP  $repo@$current (could not determine latest version)"
		continue
	fi

	# Extract major version from latest tag (e.g., v6.2.1 -> v6)
	if [[ "$latest_tag" =~ ^v([0-9]+) ]]; then
		latest="v${BASH_REMATCH[1]}"
	else
		echo "  SKIP  $repo@$current (unexpected tag format: $latest_tag)"
		continue
	fi

	if [ "$current" = "$latest" ]; then
		echo "  OK    $repo@$current (latest: $latest_tag)"
	else
		echo "  UPDATE $repo@$current -> $latest (latest: $latest_tag)"

		# Replace all occurrences across workflow files
		for file in "$WORKFLOWS_DIR"/*.yml; do
			sed -i '' "s|$repo@$current|$repo@$latest|g" "$file"
		done

		updated=$((updated + 1))
	fi
done

echo ""
if [ "$updated" -gt 0 ]; then
	echo "Updated $updated action(s)."
	echo "Run 'make test/snapshot/update' to regenerate the snapshot."
else
	echo "All actions are up to date."
fi
