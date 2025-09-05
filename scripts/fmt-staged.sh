#!/usr/bin/env bash
set -euo pipefail

fmt_bin=$(command -v gofumpt || command -v gofmt || true)
[ -z "${fmt_bin:-}" ] && { echo "No gofumpt/gofmt found" >&2; exit 1; }

files=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '\.go$' || true)
[ -z "$files" ] && { echo "No staged Go files" >&2; exit 0; }

echo "Formatting with $(basename "$fmt_bin")" >&2
echo "$files" | xargs "$fmt_bin" -w
echo "$files" | xargs git add --
