# Public Boundary

`open-ralphglasses` is a clean public seed, not a mirror of a private
operations repository. It is meant to show the architecture patterns behind a
multi-provider agent control plane without carrying private machine state,
credentials, tenants, local automation, or historical repo data.

## What Was Kept

- Provider catalog normalization across Codex, Claude, and Gemini-style CLIs.
- Budget estimation from explicit token counts and caller-provided limits.
- Review-only launch plans that show command arguments without starting child
  processes.
- Review-only loop plans with verification gates and stop conditions.
- Hook-gate decisions that classify proposed tool actions without executing
  hooks.
- A no-shell process runner with timeout and capped output for explicit local
  commands.
- Provider-neutral session planning, transcript inspection, and replay helpers.
- A small MCP-style manifest and in-process adapter with stable public tool
  names.
- Workspace discovery that requires explicit `.open-ralphrc` opt-in files.
- Deterministic managed-worktree path planning.

## What Was Removed

- Private provider credentials, token pools, account routing, and live billing
  ledgers.
- Machine-specific launchers, desktop/session control, browser state, and host
  names.
- Tenant data, private job-search state, local databases, vault exports, and
  personal documents.
- Internal operational workflows, private recovery loops, and unattended fleet
  automation.
- Full repository history from the private codebase.

## Why The Boundary Exists

The public project focuses on inspectable engineering primitives:

- typed command planning before execution,
- explicit permission and review boundaries,
- deterministic state files,
- small package-level contracts,
- testable adapters instead of live provider dependencies.

Those primitives are the reusable part. The private system adds deployment,
secrets, local workstation integration, and tenant-specific operations on top.

## Publication Checks

Before each public push, run:

```bash
make smoke
gitleaks detect --source . --redact
git status --short --branch
```

The CI workflow runs the same smoke path with full Git history checkout and a
gitleaks scan.
