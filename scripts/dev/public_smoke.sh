#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${repo_root}"

export GOWORK=off

go test ./...
go run . providers >/dev/null
go run . doctor >/dev/null
go run . mcp manifest >/dev/null
go run . worktree path --repo example --label smoke >/dev/null
go run . repos scan --root . >/dev/null

if command -v gitleaks >/dev/null 2>&1; then
  gitleaks detect --source . --redact
fi

echo "public smoke passed"
