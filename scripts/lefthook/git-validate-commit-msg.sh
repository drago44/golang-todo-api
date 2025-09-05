#!/usr/bin/env bash
set -euo pipefail

# Git passes the commit message file path as the first argument
msg_file="$1"
first_line=$(sed -n '1p' "$msg_file" | tr -d '\r')

# Allow merge commits, fixups, and reverts
if [[ "$first_line" =~ ^(Merge|Revert|fixup!|squash!) ]]; then
  exit 0
fi

# Extract ticket from branch if present: BE-/FS-/FE- followed by digits (force uppercase)
branch_name=$(git rev-parse --abbrev-ref HEAD)
branch_upper=$(printf "%s" "$branch_name" | tr '[:lower:]' '[:upper:]')
ticket=""
if [[ "$branch_upper" =~ ^((BE|FS|FE)-[0-9]+) ]]; then
  ticket="${BASH_REMATCH[1]}"
fi

# If message already starts with a ticket, accept
if [[ "$first_line" =~ ^(BE|FS|FE)-[0-9]+\b ]]; then
  exit 0
fi

# If we have a ticket from branch, auto-prefix it to the first line, preserving existing text
if [[ -n "$ticket" ]]; then
  rest=$(sed '1d' "$msg_file")
  printf "%s %s\n" "$ticket" "$first_line" > "$msg_file"
  if [[ -n "$rest" ]]; then
    printf "%s\n" "$rest" >> "$msg_file"
  fi
  exit 0
fi

# Otherwise, do not enforce when no ticket is derivable (e.g., dev branch)
exit 0
