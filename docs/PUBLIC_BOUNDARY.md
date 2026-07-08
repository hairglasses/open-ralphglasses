# Public Boundary

`open-ralphglasses` is a public seed for reviewable agent planning. It keeps
the small pieces that are useful without any account-specific state or host
assumptions.

## Included

- Provider catalog checks for Codex, Claude, and Gemini.
- Token-cost estimation from explicit token counts.
- Hook-gate review for proposed tool actions.
- Review-only launch plans that return command arguments without executing
  child processes.
- Review-only loop plans with bounded iterations and verification commands.
- MCP-style tool names and in-process dispatch for the public commands.

## Excluded

- Credentials, tokens, browser session data, and account exports.
- Operator-specific paths, hostnames, caches, or machine bootstrap code.
- Live provider process supervision or session recovery. The production
  crash/interrupt-recovery primitives for this are being extracted separately
  as `durable-recovery` (in progress).
- Repo discovery, worktree control, transcript parsing, or shell execution.
  The worktree-lifecycle piece of this is already public and standalone as
  `worktree-isolate`.
- An interactive TUI/session viewer, and a real hook-execution runtime — this
  repo's `internal/hookgate` only reviews and classifies proposed actions
  (allow/warn/block); turning that into enforced, running hooks is a
  separate, in-progress project (`hook-runner`).
- A fuller MCP transport (stdio/HTTP servers, auth, rate limiting, routing)
  instead of the in-process adapter here. `mcp-gateway` is the
  production-shaped gateway this seed's MCP-style surface points toward.
- Non-public research pipelines and generated provider overlays.

See the README's "Extracted Components" section for the full list of
standalone projects this system's ideas have grown into. That section
deliberately omits hosting-organization links: this repo's boundary check
(`scripts/dev/public_smoke.sh`) forbids the parent workspace's brand name
appearing anywhere in this tree, and that check is authoritative — do not
weaken it to add convenience links without an explicit, deliberate decision
outside of routine doc edits.

## Publication Checks

Before expanding the public surface, run:

```bash
go test ./...
bash scripts/dev/public_smoke.sh
gitleaks detect --source . --no-git --redact
```

New features should preserve the same rule: inputs are explicit, outputs are
reviewable JSON or plain text, and commands do not execute external providers
unless a future feature deliberately documents and tests that behavior.
