# open-ralphglasses

Public, sanitized seed of a multi-provider agent control plane.

This repository keeps the public-safe core ideas from a larger private
ralphglasses system:

- normalize multiple agent providers behind one small catalog
- estimate provider token cost and simple budget headroom
- evaluate public hook-gate decisions for proposed tool actions
- build review-only provider launch plans without starting child processes
- build review-only loop plans with verification gates and stop conditions
- run explicit no-shell local processes with timeout and capped output
- validate and record planned provider sessions
- keep local session state in simple JSONL files
- publish a tiny MCP-style command manifest
- scan Git workspaces for explicit public opt-in markers
- compute deterministic managed worktree paths

It intentionally does not include private automation, credential routing,
machine-specific launchers, private docs, tenant state, secret stores, local
browser state, or internal operational workflows.

## Install

```bash
go install github.com/hairglasses/open-ralphglasses@latest
```

From source:

```bash
git clone https://github.com/hairglasses/open-ralphglasses.git
cd open-ralphglasses
go test ./...
go run . providers
```

## Usage

```bash
# Check which provider CLIs are installed on this machine.
go run . doctor

# List the provider catalog.
go run . providers

# Estimate token cost and optional budget headroom.
go run . budget estimate --provider codex --input-tokens 1000 --output-tokens 500 --budget 1 --spent 0.25

# Evaluate a proposed hook action without executing hooks.
go run . hook check --event PreToolUse --tool Bash --input "git status --short"

# Build an inspectable provider command plan without executing it.
go run . launch plan --provider codex --repo . --prompt "Summarize this repository" --permission-mode read-only

# Build an inspectable iterative work plan without executing it.
go run . loop plan --repo . --goal "Improve tests" --provider codex --verify "go test ./..."

# Run one explicit process without a shell.
go run . process run --repo . --timeout-seconds 10 -- go version

# Record a planned session without launching a child process.
go run . session start --provider codex --repo . --prompt "Summarize this repository"

# Read planned sessions from .open-ralph/sessions.jsonl.
go run . session list

# Emit the public MCP-style tool manifest.
go run . mcp manifest

# Scan a workspace for Git repos and .open-ralphrc opt-in files.
go run . repos scan --root . --depth 3

# Compute a deterministic managed worktree path.
go run . worktree path --repo example-service --label refactor-api
```

## Public Boundary

The first public version is deliberately small. It is a clean-room public subset,
not a full mirror of the private repository history. The code and docs are
designed to be inspectable, example-driven, and free of private references.

## Development

```bash
make smoke
gitleaks detect --source . --no-git --redact
```
