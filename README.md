# open-ralphglasses

[![CI](https://github.com/hairglasses/open-ralphglasses/actions/workflows/ci.yml/badge.svg)](https://github.com/hairglasses/open-ralphglasses/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Small Go CLI for planning multi-provider agent work.

This is a sanitized public seed of a broader control-plane direction. It keeps the reusable planning, policy, and MCP-style contract surfaces while excluding account data, customer data, host details, and non-public research pipelines.

The project focuses on reviewable planning primitives:

- list supported agent provider CLIs
- estimate token cost from explicit counts
- classify proposed hook actions as allow, warn, or block
- build provider command plans without executing them
- build bounded loop plans with verification gates
- expose the same core operations through a small MCP-style adapter

## Install

Stable release:

```bash
go install github.com/hairglasses/open-ralphglasses@v0.1.0
```

Latest main:

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

# Emit the MCP-style tool manifest.
go run . mcp manifest

# Call an MCP-style tool in process.
go run . mcp call open_ralph_budget_estimate --param provider=codex --param input_tokens=1000 --param output_tokens=500
```

## Demo Output

`open-ralphglasses` returns reviewable plans and checks; it does not execute
provider agents or shell hooks.

```text
$ open-ralphglasses doctor
provider  installed  command
codex     true       codex
claude    true       claude
gemini    true       gemini
```

```json
{
  "provider": "codex",
  "command": "codex",
  "args": [
    "exec",
    "--model",
    "gpt-5",
    "--json",
    "--sandbox",
    "read-only",
    "Summarize this repository"
  ],
  "repo_path": "/path/to/open-ralphglasses",
  "permission_mode": "read-only",
  "review_required": true,
  "safety_notes": [
    "command is not executed by open-ralphglasses",
    "review args and environment policy before wiring a process runner"
  ],
  "execution_status": "planned_only"
}
```

```json
{
  "goal": "Improve tests",
  "repo_path": "/path/to/open-ralphglasses",
  "provider": "codex",
  "max_iterations": 3,
  "verify_commands": [
    "go test ./..."
  ],
  "review_required": true,
  "execution_status": "planned_only",
  "steps": [
    {
      "order": 1,
      "name": "inspect",
      "description": "read repo state, task context, and current diffs"
    },
    {
      "order": 2,
      "name": "plan_slice",
      "description": "choose one bounded implementation slice"
    },
    {
      "order": 3,
      "name": "implement",
      "description": "apply the smallest changes needed for the slice"
    },
    {
      "order": 4,
      "name": "verify",
      "description": "run the configured verification command and capture output"
    }
  ],
  "stop_conditions": [
    "goal is satisfied",
    "verification fails and requires human review",
    "max_iterations is reached"
  ]
}
```

## Development

```bash
make smoke
gitleaks detect --source . --no-git --redact
```

## Public Surface

- `docs/EXAMPLES.md` shows public-safe output shapes.
- `docs/PUBLIC_BOUNDARY.md` documents what is intentionally included and
  excluded from this reduced public seed.
