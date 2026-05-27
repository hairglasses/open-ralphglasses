# MCP-Style Manifest

This repository ships a stable tool manifest and an in-process call adapter so
future stdio or HTTP adapters can preserve the command names.

Current tools:

- `open_ralph_doctor`
- `open_ralph_tool_manifest`
- `open_ralph_provider_list`
- `open_ralph_budget_estimate`
- `open_ralph_hook_check`
- `open_ralph_launch_plan`
- `open_ralph_loop_plan`

Run:

```bash
go run . mcp manifest
go run . mcp call open_ralph_provider_list
```

See `docs/EXAMPLES.md` for safe output examples and
`docs/PUBLIC_BOUNDARY.md` for the current public cut line.
