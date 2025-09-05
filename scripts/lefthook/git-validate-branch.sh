#!/usr/bin/env bash
set -euo pipefail

branch_name=$(git rev-parse --abbrev-ref HEAD)

# Allow detached HEAD (e.g., CI) without failing
if [[ "$branch_name" == "HEAD" ]]; then
  exit 0
fi

if [[ "$branch_name" == "dev" || "$branch_name" =~ ^(BE|FS|FE)-[0-9]+ ]]; then
  exit 0
fi

echo "Branch name '$branch_name' is invalid." 1>&2
echo "Expected 'dev' or start with BE-, FS-, FE- followed by digits, e.g.:" 1>&2
echo "  BE-1234-some-feature" 1>&2
echo "  FS-42-bugfix" 1>&2
echo "  FE-7-ui-tweak" 1>&2
exit 1
