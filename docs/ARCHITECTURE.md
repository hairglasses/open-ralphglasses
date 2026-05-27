# Architecture

`open-ralphglasses` is split into small packages with explicit public contracts.

| Package | Purpose |
|---|---|
| `internal/provider` | Provider catalog and validation. |
| `internal/budget` | Public cost estimates and budget headroom checks. |
| `internal/hookgate` | Public allow/warn/block decisions for proposed hook actions. |
| `internal/launchplan` | Review-only provider CLI command planning. |
| `internal/loopplan` | Review-only iterative work plans with verification gates. |
| `internal/mcpadapter` | In-process dispatch for public MCP-style tool names. |
| `internal/processrun` | Explicit no-shell process execution with timeout and capped output. |
| `internal/session` | Durable JSONL session planning ledger. |
| `internal/events` | Bounded in-memory event history for adapters. |
| `internal/discovery` | Public workspace scan over Git repos and `.open-ralphrc` opt-in files. |
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
3. Command packages normalize repo, prompt, hook, loop, process, discovery, or
   path metadata.
4. `internal/session.Store` appends session plans under `.open-ralph/` when the
   session command is used.
5. Future TUI, MCP, or HTTP adapters can read the same public package outputs.

## Extension Points

- Add a provider by extending `provider.Catalog`.
- Add provider pricing by passing explicit rates into `internal/budget` instead
  of treating example rates as current billing truth.
- Add process execution only above `internal/launchplan`, after deciding how the
  caller will review environment and permission policy.
- Add a command by wiring a package function in `internal/cli`.
- Add MCP transport by adapting `mcpmanifest.Manifest`.
- Add repo metadata by teaching `internal/discovery` about another explicit,
  public-safe marker file.
- Add real launches by building an explicit process package above
  `session.New`, with tests for command construction and cancellation.
