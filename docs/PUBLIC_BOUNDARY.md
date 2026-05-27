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
- Live provider process supervision or session recovery.
- Repo discovery, worktree control, transcript parsing, or shell execution.
- Non-public research pipelines and generated provider overlays.

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
