# Getting Started

## Requirements

- Go 1.23 or newer
- Optional provider CLIs on `PATH`: `codex`, `claude`, or `gemini`

## First Run

```bash
go test ./...
go run . doctor
go run . providers
go run . budget estimate --provider codex --input-tokens 1000 --output-tokens 500
go run . hook check --event PreToolUse --tool Bash --input "git status --short"
go run . launch plan --provider codex --repo . --prompt "Inspect this repository" --permission-mode read-only
go run . loop plan --repo . --goal "Improve tests" --provider codex --verify "go test ./..."
go run . mcp manifest
go run . mcp call open_ralph_provider_list
```

## Next Steps

- Use `go run . budget estimate` for a JSON-readable cost estimate.
- Use `go run . hook check` to review whether a proposed tool action should be
  allowed, warned, or blocked.
- Use `go run . launch plan` to inspect provider-specific CLI args before
  executing anything elsewhere.
- Use `go run . loop plan` to describe bounded implementation iterations,
  verification gates, and stop conditions.
- Read `docs/EXAMPLES.md` for public-safe output examples.
- Read `docs/PUBLIC_BOUNDARY.md` before adding new surfaces.
- Read `docs/ARCHITECTURE.md` before adding new planning or transport code.
