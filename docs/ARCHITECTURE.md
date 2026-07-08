# Architecture

`open-ralphglasses` is split into small packages with explicit contracts.

| Package | Purpose |
|---|---|
| `internal/provider` | Provider catalog and validation. |
| `internal/budget` | Cost estimates and budget headroom checks. |
| `internal/hookgate` | Allow, warn, and block decisions for proposed hook actions. |
| `internal/launchplan` | Review-only provider CLI command planning. |
| `internal/loopplan` | Review-only iterative work plans with verification gates. |
| `internal/mcpadapter` | In-process dispatch for MCP-style tool names. |
| `internal/mcpmanifest` | Machine-readable command manifest. |
| `internal/cli` | Thin command routing over the packages above. |

## Data Flow

1. The CLI parses command flags.
2. Package-level validation normalizes provider ids, repo paths, prompts, and
   option values.
3. Planning packages return JSON-friendly structs for command plans, loop
   plans, cost estimates, and hook decisions.
4. The MCP-style adapter dispatches to the same packages so command behavior
   stays consistent across entry points.

## Extension Points

- Add a provider by extending `provider.Catalog`.
- Add provider pricing by passing explicit rates into `budget.Estimate`.
- Add a planning command by wiring a package function in `internal/cli`.
- Add transport support by adapting `mcpmanifest.Manifest` and
  `mcpadapter.Call`.

## Relationship To The Broader System

This repo intentionally stays small: a reviewable planning and policy layer
with no live execution. Several of the harder, production-shaped problems
this planning layer assumes have been solved are developed and published as
their own standalone repos rather than folded back in here — see the
README's "Extracted Components" section for the full list (`mcpkit`,
`agentloop`, `provider-quota`, `leasequeue`, `mcp-gateway`,
`worktree-isolate`, `syspressure`, `agent-consensus`, `provider-shim`,
`browserctl`, `agent-patterns`, and the in-progress `durable-recovery` and
`hook-runner`). Treat `open-ralphglasses` as the conceptual seed and those
repos as the components that grew out of it, not as a single monolith.
