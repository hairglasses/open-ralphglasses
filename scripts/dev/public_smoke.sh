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

for encoded in \
  czJm \
  a2lh \
  c2VjcmV0c3R1ZGlvcw== \
  cnVubXlsaWZl \
  am9iYg== \
  Z21haWw= \
  bGlua2VkaW4= \
  cmVzdW1l \
  dGVuYW50 \
  dmF1bHQ= \
  YnJvd3NlciBzdGF0ZQ== \
  cHJpdmF0ZSByYWxwaGdsYXNzZXM= \
  cHJpdmF0ZSBwcm9qZWN0 \
  cHJpdmF0ZSBzeXN0ZW0= \
  bGFyZ2VyIHByaXZhdGU= \
  cGVyc29uYWw= \
  bWFjaGluZS1zcGVjaWZpYw== \
  aW50ZXJuYWwgb3BlcmF0aW9uYWw= \
  aGFpcmdsYXNzZXMtc3R1ZGlv \
  YXJjaGdsYXNzZXM=
do
  marker="$(printf '%s' "${encoded}" | base64 -d)"
  if rg -n -i --fixed-strings --hidden --glob '!/.git/**' -- "${marker}" .; then
    echo "public marker guard failed" >&2
    exit 1
  fi
done

if command -v gitleaks >/dev/null 2>&1; then
  gitleaks detect --source . --redact
fi

echo "public smoke passed"
