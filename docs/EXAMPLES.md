# Example Outputs

These examples are safe public fixtures for the reduced CLI surface. They are
intended to show output shape, not current provider pricing authority.

## Provider Catalog

```text
$ go run . providers
ID      DISPLAY           COMMAND  DEFAULT MODEL  NOTES
codex   OpenAI Codex CLI  codex    gpt-5          coding-oriented agent runtime
claude  Claude Code       claude   sonnet         interactive coding agent runtime
gemini  Gemini CLI        gemini   gemini-pro     general-purpose agent runtime
```

## Budget Estimate

```json
{
  "estimate": {
    "provider": "codex",
    "input_tokens": 1000,
    "output_tokens": 500,
    "estimated_usd": 0.01,
    "baseline_usd": 0.0105,
    "efficiency_pct": 95.24,
    "rate_input_per_1m_usd": 2.5,
    "rate_output_per_1m_usd": 15,
    "rate_source": "example"
  },
  "status": {
    "spent_usd": 0.25,
    "budget_usd": 1,
    "remaining_usd": 0.75,
    "percent_used": 25,
    "headroom_pct": 90,
    "should_stop": false,
    "budget_disabled": false
  }
}
```

## Hook Check

```json
{
  "event": "PreToolUse",
  "tool": "Bash",
  "verdict": "allow",
  "review_required": false,
  "reasons": [
    "no public hook-gate rules matched"
  ],
  "recommendations": [
    "review the decision before wiring it into an automatic hook runner"
  ]
}
```

## MCP-Style Manifest

```json
[
  {
    "name": "open_ralph_tool_manifest",
    "description": "Return the public MCP-style tool manifest",
    "read_only": true
  },
  {
    "name": "open_ralph_doctor",
    "description": "Check provider executable availability",
    "read_only": true
  },
  {
    "name": "open_ralph_provider_list",
    "description": "List configured public providers",
    "read_only": true
  },
  {
    "name": "open_ralph_budget_estimate",
    "description": "Estimate provider token cost and evaluate optional budget headroom",
    "read_only": true,
    "inputs": [
      "provider",
      "input_tokens",
      "output_tokens",
      "budget",
      "spent"
    ]
  }
]
```
