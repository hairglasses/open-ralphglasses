# Portfolio Notes

This repository is intended to be readable by interviewers and profile viewers.
It demonstrates the public-safe slice of a larger agent-control-plane project.

## What To Look At

- `internal/provider`: provider-neutral catalog and validation.
- `internal/launchplan`: command planning that keeps execution reviewable.
- `internal/loopplan`: bounded iteration plans with verification gates.
- `internal/hookgate`: allow, warn, and block decisions for proposed tool use.
- `internal/processrun`: explicit no-shell execution with timeout and capped
  output.
- `internal/mcpadapter` and `internal/mcpmanifest`: stable MCP-style public
  tool contracts.
- `internal/sessionlog`: provider-neutral transcript summaries and replay text.

## Resume Claim Mapping

Use these as evidence-backed wording anchors:

- Built a Go control-plane seed for multi-provider AI agent workflows.
- Designed review-first command planning for provider launches and iterative
  loops.
- Implemented budget and safety primitives for agent sessions without requiring
  live credentials.
- Added MCP-style tool contracts and an in-process adapter for testable tool
  dispatch.
- Preserved a strict public/private boundary while extracting reusable
  architecture from a larger production system.

## Suggested Demo Flow

```bash
go run . providers
go run . budget estimate --provider codex --input-tokens 1000 --output-tokens 500 --budget 1
go run . launch plan --provider codex --repo . --prompt "Summarize this repository" --permission-mode read-only
go run . loop plan --repo . --goal "Improve tests" --provider codex --verify "go test ./..."
go run . mcp manifest
make smoke
```

The demo does not require API keys or provider accounts.
