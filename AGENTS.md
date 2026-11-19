# open-ralphglasses

> Canonical instructions: AGENTS.md

Public Go CLI (`github.com/hairglasses/open-ralphglasses`) for planning multi-provider agent work: provider discovery, token-budget estimation, review-only hook-action classification, launch/loop plan building, and an MCP-style tool manifest/call adapter.

## Layout

- `main.go`, `internal/cli` — CLI entrypoint and subcommand wiring.
- `internal/provider` — provider catalog and doctor checks.
- `internal/budget` — token-cost estimation.
- `internal/hookgate` — review-only hook-action classification (`allow`/`warn`/`block`).
- `internal/launchplan`, `internal/loopplan` — dry-run launch and bounded loop plan builders.
- `internal/mcpadapter`, `internal/mcpmanifest` — MCP-style tool manifest and in-process call adapter.
- `docs/` — `ARCHITECTURE.md`, `EXAMPLES.md`, `GETTING_STARTED.md`, `MCP_TOOLS.md`, `PORTFOLIO_PROOF.md`, `PUBLIC_BOUNDARY.md`.

## Boundaries & Gates

- **Public boundary**: Read `docs/PUBLIC_BOUNDARY.md` before adding surfaces. `scripts/dev/public_smoke.sh` checks that no non-public markers or org links leak into tracked files.
- **Plans only**: Keep `launch plan`, `loop plan`, and `hook check` strictly non-executing.
- **Verify before committing**:
  - `make test` (`GOWORK=off go test ./...`)
  - `make smoke` (`bash scripts/dev/public_smoke.sh`)
  - `gitleaks detect --source . --no-git --redact`

## Landing

- Land verified changes via `branch -> PR -> claude-harness/bin/hg-land --repo <checkout>` (never `git push origin main` directly; never use `git add -A`).

<!-- BEGIN GENERATED LANDING POLICY -->
Use `hg-land --repo <checkout>` (`--no-verify` while the fleet gate switch
`~/.config/hg/gates-disabled` or `HG_GATES_DISABLED=1` is active; `--local-gate <log>`
under `HG_GATES_DISABLED=0`, or `hg-gate-then-land --require-gate`) to land through a
pull request. Owned by `claude-harness/.agents/policies/landing.md`.

- While fleet gates are disabled (operator 2026-11-08), `hg-land` lands without a
  local gate log; when gated, bind the log to `head=<full SHA>` with final `rc=0`.
- Observe the repository's current remote requirements. Required checks and
  branch protections remain binding; hosted CI is not assumed required everywhere.
- Never bypass a failed gate, force-push, or replace unrelated work. Do not use
  direct pushes to main as a substitute for the configured landing route.
- Optoinvest repositories remain local-only: no remote publication.
- Prove landing with the tool receipt and exact remote main SHA. Source updates
  and runtime/deployment verification are separate claims.
<!-- END GENERATED LANDING POLICY -->
