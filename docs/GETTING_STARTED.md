# Getting Started

## Requirements

- Go 1.23 or newer
- Optional provider CLIs on `PATH`: `codex`, `claude`, or `gemini`

## First Run

```bash
go test ./...
go run . doctor
go run . budget estimate --provider codex --input-tokens 1000 --output-tokens 500
go run . launch plan --provider codex --repo . --prompt "Inspect this repository" --permission-mode read-only
go run . session start --provider codex --repo . --prompt "Inspect this repository"
go run . session list
go run . repos scan --root . --depth 3
```

The session command writes `.open-ralph/sessions.jsonl`. That path is ignored by
git so local prompts and repo paths stay local.

## Next Steps

- Use `go run . mcp manifest` to inspect the public command surface.
- Use `go run . budget estimate` when you want a local, JSON-readable cost
  estimate before wiring real process execution.
- Use `go run . launch plan` to review provider-specific CLI args before adding
  a real process runner.
- Add `.open-ralphrc` to example repositories you want surfaced as enabled by
  `go run . repos scan`.
- Use `go run . worktree path --repo my-service --label add-tests` to preview
  where a managed worktree would live.
- Read `docs/ARCHITECTURE.md` before adding launch or transport code.
