# open-ralphglasses

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
