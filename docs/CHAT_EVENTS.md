# Chat Events

`open-ralphglasses` stores provider-neutral transcript artifacts under
`.open-ralph/transcripts/<session-id>.json`.

The transcript schema is intentionally separate from the planned-session ledger
in `.open-ralph/sessions.jsonl`. Planned sessions describe launch intent.
Transcript artifacts describe events captured by a future runner or imported by
tests and tools.

## Event Kinds

- `operator_message` - prompt text supplied by the operator.
- `start` - provider turn start metadata.
- `delta` - additive assistant text or thinking text.
- `tool_use_start` - structural tool-call boundary.
- `tool_use_end` - structural tool-result boundary.
- `usage` - token accounting update.
- `end` - provider turn ended.
- `error` - stream or provider error.

Structural events must not be dropped during replay. Additive `delta` and
`usage` events may be merged by a future streaming layer.

## Commands

```bash
go run . session inspect --id sess-example
go run . session analyze --id sess-example
go run . session replay-text --id sess-example
```

These commands only read local transcript artifacts. They do not launch
providers, resume sessions, read credentials, or inspect browser state.
