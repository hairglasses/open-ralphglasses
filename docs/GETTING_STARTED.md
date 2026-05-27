# Getting Started

## Requirements

- Go 1.23 or newer
- Optional provider CLIs on `PATH`: `codex`, `claude`, or `gemini`

## First Run

```bash
go test ./...
go run . doctor
go run . session start --provider codex --repo . --prompt "Inspect this repository"
go run . session list
```

The session command writes `.open-ralph/sessions.jsonl`. That path is ignored by
git so local prompts and repo paths stay local.

## Next Steps

- Use `go run . mcp manifest` to inspect the public command surface.
- Use `go run . worktree path --repo my-service --label add-tests` to preview
  where a managed worktree would live.
- Read `docs/ARCHITECTURE.md` before adding launch or transport code.

