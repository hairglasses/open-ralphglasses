# open-ralphglasses — Agent Instructions

> Canonical instructions: AGENTS.md

This repo uses AGENTS.md as the canonical instruction file.

## What this is

A small public Go CLI (`github.com/hairglasses/open-ralphglasses`) for
planning multi-provider agent work: provider discovery, token-budget
estimation, hook-action classification, provider command-plan building,
bounded loop planning, and an MCP-style tool manifest/call adapter. It
reviews and plans — it does not execute provider agents or run shell hooks.

**This repo is a publish mirror.** Canonical development of the system it
mirrors happens in the private `ralphglasses` repo. `open-ralphglasses` is a
deliberately reduced, sanitized public seed: it keeps the reusable
planning/policy/MCP-contract surfaces and excludes account data, customer
data, host details, and non-public research pipelines. Several of its
underlying ideas have since grown into their own standalone, independently
published Go projects (see the README's "Extracted Components" table) —
those are separate repos, not renames of this one.

## Layout

- `main.go` — CLI entrypoint.
- `internal/cli` — command wiring.
- `internal/provider` — provider catalog + doctor checks.
- `internal/budget` — token-cost estimation.
- `internal/hookgate` — review-only hook-action classification (allow/warn/block).
- `internal/launchplan` — provider command-plan construction (no execution).
- `internal/loopplan` — bounded iterative-work plan construction.
- `internal/mcpadapter`, `internal/mcpmanifest` — MCP-style tool manifest and in-process call adapter.
- `docs/` — `GETTING_STARTED.md`, `EXAMPLES.md`, `ARCHITECTURE.md`, `PORTFOLIO_PROOF.md`, `PUBLIC_BOUNDARY.md`, `MCP_TOOLS.md`.
- `scripts/dev/public_smoke.sh` — the boundary check gate (see below).

## Conventions

- **Public boundary is a hard gate.** Before adding new planning surfaces,
  read `docs/PUBLIC_BOUNDARY.md`. `scripts/dev/public_smoke.sh` treats the
  parent workspace's brand name as a forbidden marker and fails the check if
  it leaks in — that's intentional, not a bug to work around.
- **Plans only, never execution.** Every planning surface (`launch plan`,
  `loop plan`, `hook check`) must remain inspectable-output-only. Don't wire
  these to actually spawn provider CLIs or run hooks; that would break the
  repo's reason for being public.
- **Don't add hosting-org links for extracted components.** The README
  intentionally doesn't link out to a specific GitHub org for the extracted
  standalone projects — the boundary check enforces this. Reference them by
  name only.
- Before committing: `make test`, `make smoke`, and
  `gitleaks detect --source . --no-git --redact`.

## Workspace context

See `/home/hg/hairglasses-studio/CLAUDE.md` for shared multi-repo
conventions. The private repo this mirrors is
`/home/hg/hairglasses-studio/ralphglasses`. Do not sync private-repo
internals here wholesale — only intentionally reduced, sanitized slices,
reviewed against `docs/PUBLIC_BOUNDARY.md`.
