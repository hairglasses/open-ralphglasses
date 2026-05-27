# Architecture

`open-ralphglasses` is split into small packages with explicit public contracts.

| Package | Purpose |
|---|---|
| `internal/provider` | Provider catalog and validation. |
| `internal/session` | Durable JSONL session planning ledger. |
| `internal/events` | Bounded in-memory event history for adapters. |
| `internal/worktree` | Deterministic managed worktree path planning. |
| `internal/mcpmanifest` | Machine-readable public command manifest. |
| `internal/cli` | Thin command routing over the packages above. |

The current version records planned sessions rather than launching provider
processes. That is intentional: process execution, credential policy, remote
workers, browser automation, and tenant state need separate public design before
they belong in this repository.

## Data Flow

1. CLI parses public command flags.
2. `internal/provider` validates the provider id.
3. `internal/session` normalizes repo and prompt metadata.
4. `internal/session.Store` appends a JSONL record under `.open-ralph/`.
5. Future TUI, MCP, or HTTP adapters can read the same session ledger.

## Extension Points

- Add a provider by extending `provider.Catalog`.
- Add a command by wiring a package function in `internal/cli`.
- Add MCP transport by adapting `mcpmanifest.Manifest`.
- Add real launches by building an explicit process package above
  `session.New`, with tests for command construction and cancellation.

