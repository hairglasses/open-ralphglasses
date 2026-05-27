# Public MCP-Style Manifest

This repository does not ship a full MCP transport yet. It does ship a stable
tool manifest and an in-process call adapter so a future stdio or HTTP adapter
can preserve the command names.

Current public tools:

- `open_ralph_doctor`
- `open_ralph_tool_manifest`
- `open_ralph_provider_list`
- `open_ralph_process_run`
- `open_ralph_budget_estimate`
- `open_ralph_hook_check`
- `open_ralph_launch_plan`
- `open_ralph_loop_plan`
- `open_ralph_session_plan`
- `open_ralph_session_list`
- `open_ralph_repo_scan`
- `open_ralph_worktree_path`

Run:

```bash
go run . mcp manifest
go run . mcp call open_ralph_provider_list
```
