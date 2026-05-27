# open-ralphglasses

Public, sanitized seed of a multi-provider agent control plane.

This repository keeps the public-safe core ideas from a larger private
ralphglasses system:

- normalize multiple agent providers behind one small catalog
- estimate provider token cost and simple budget headroom
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
