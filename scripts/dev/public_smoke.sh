#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${repo_root}"

export GOWORK=off

go test ./...
go run . budget estimate --provider codex --input-tokens 1000 --output-tokens 500 >/dev/null
go run . hook check --event PreToolUse --tool Bash --input "git status --short" >/dev/null
go run . launch plan --provider codex --repo . --prompt "Summarize this repository" --permission-mode read-only >/dev/null
go run . loop plan --repo . --goal "Improve docs" --provider codex --verify "go test ./..." >/dev/null
go run . providers >/dev/null
go run . doctor >/dev/null
go run . mcp manifest >/dev/null
go run . process run --repo . --timeout-seconds 10 -- go version >/dev/null
go run . worktree path --repo example --label smoke >/dev/null
go run . repos scan --root . >/dev/null

if command -v gitleaks >/dev/null 2>&1; then
  gitleaks detect --source . --redact
fi

echo "public smoke passed"
