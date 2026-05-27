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
go run . mcp call open_ralph_budget_estimate --param provider=codex --param input_tokens=1000 --param output_tokens=500 >/dev/null
go run . process run --repo . --timeout-seconds 10 -- go version >/dev/null
tmp_root="$(mktemp -d)"
trap 'rm -rf "${tmp_root}"' EXIT
mkdir -p "${tmp_root}/.open-ralph/transcripts"
cat >"${tmp_root}/.open-ralph/transcripts/sess-smoke.json" <<'JSON'
{
  "schema_version": 1,
  "snapshot": {
    "id": "sess-smoke",
    "provider": "codex",
    "provider_session_id": "provider-smoke"
  },
  "transcript": [
    {"kind": "operator_message", "session_id": "sess-smoke", "text": "Inspect this repository", "at": "2026-05-27T00:00:00Z"},
    {"kind": "delta", "session_id": "sess-smoke", "channel": "text", "text": "Done", "at": "2026-05-27T00:00:01Z"},
    {"kind": "end", "session_id": "sess-smoke", "stop_reason": "end_turn", "at": "2026-05-27T00:00:02Z"}
  ]
}
JSON
go run . session analyze --root "${tmp_root}" --id sess-smoke >/dev/null
go run . session replay-text --root "${tmp_root}" --id sess-smoke >/dev/null
go run . mcp call open_ralph_session_analyze --param root="${tmp_root}" --param id=sess-smoke >/dev/null
go run . worktree path --repo example --label smoke >/dev/null
go run . repos scan --root . >/dev/null

if command -v gitleaks >/dev/null 2>&1; then
  gitleaks detect --source . --redact
fi

echo "public smoke passed"
