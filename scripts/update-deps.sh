#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEMPLATE_DIR="$REPO_ROOT/{{ .ProjectKebab }}"
TMP_DIR="$(mktemp -d)"

trap 'rm -rf "$TMP_DIR"' EXIT

# Render the scaffold using the test preset
echo "Scaffolding to $TMP_DIR..."
scaffold new --output-dir="$TMP_DIR" --preset="test" --no-prompt "$REPO_ROOT"

# The test preset uses Project: "cli-test"
PROJECT_DIR="$TMP_DIR/cli-test"

if [ ! -d "$PROJECT_DIR" ]; then
	echo "error: scaffolded project not found at $PROJECT_DIR"
	exit 1
fi

# Run go mod tidy in the rendered project
echo "Running go mod tidy..."
cd "$PROJECT_DIR"
go mod tidy

# Copy go.sum back directly (no template variables)
cp "$PROJECT_DIR/go.sum" "$TEMPLATE_DIR/go.sum"

# Copy go.mod back, restoring the template module line
GOMOD_MODULE="github.com/username/project"
sed "1s|module ${GOMOD_MODULE}|module {{ .Scaffold.gomod }}|" "$PROJECT_DIR/go.mod" > "$TEMPLATE_DIR/go.mod"

echo "Updated go.mod and go.sum in template directory."
echo ""
echo "Run 'make test/snapshot/update' to regenerate the snapshot."
