# Lefthook setup

This project uses Lefthook to enforce checks on commit/push.

## Install

- Cross-platform (recommended):
  - `go install github.com/evilmartians/lefthook@latest`
- macOS (Homebrew):
  - `brew install lefthook`

Install git hooks:

```bash
lefthook install
```

## Hooks

### pre-commit (sequential)

1. Branch name validation

- Script: `scripts/lefthook/git-validate-branch.sh`
- Allowed: `dev` or `^(BE|FS|FE)-[0-9]+` (e.g., `BE-1234-feature`)

2. Format staged Go files

- Script: `scripts/lefthook/fmt-staged.sh`
- Uses `gofumpt` if available, otherwise `gofmt`
- Formats only staged `.go` files and re-adds them to the index
- If no staged changes remain after formatting, aborts commit to avoid empty commits

3. Vet

- Command: `make vet`

4. Tests (short)

- Command: `make test-short`
- Runs `gotestsum` if installed, otherwise falls back to `go test`

5. Lint

- Command: `make lint`
- Runs `golangci-lint` if installed, otherwise fails with install hint

### commit-msg

- Script: `scripts/lefthook/git-validate-commit-msg.sh`
- Extracts ticket prefix from the current branch (BE/FS/FE + id), uppercases it, and prepends to the first message line without removing existing text
- Skips for `Merge`, `Revert`, `fixup!`, `squash!`

### pre-push

- Command: `make test-short`

## Make targets used by hooks

- `make test`: richgo if available, fallback to `go test`
- `make test-short`: `gotestsum` short-verbose if available, fallback to `go test`
- `make fmt`: applies formatting (gofumpt if installed, else go fmt)
- `make vet`, `make lint`

## Troubleshooting

- Hooks not running in IDE: ensure `$HOME/go/bin` is in PATH, or install Lefthook via Homebrew
- To re-install hooks after config changes:

```bash
lefthook install
```

- Run hooks manually:

```bash
lefthook run pre-commit --all-files
lefthook run commit-msg /path/to/message.txt
```
