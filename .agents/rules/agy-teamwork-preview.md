---
trigger: always_on
---

# AGY 2.0 Teamwork Preview & Delegation Rules

> ⚠️ **`teamwork-preview` is NOT an agy execution mode.** Verified against the installed
> binary on 2026-08-18: `agy --help` reported
> `--mode  Set the agent execution mode for this session (accept-edits, plan)`.
> Re-verified on 2026-10-22 (agy on ag1): the same line now reads
> `--mode  Set the agent execution mode for this session (accept-edits, auto, plan)`.
> The set has grown by `auto`, but `teamwork-preview` is in neither list. Do not pass
> `--mode teamwork-preview`; it will not parse. `teamwork-preview` exists only as the name
> of this rule file and a docs directory, and the "AGY 2.0 subagent spawning spec" it
> derives from was found to be fabricated (its verification harness mocked itself and
> never executed agy; recorded 2026-08-18, not re-verified).
>
> The delegation and write-boundary guidance below is still useful as *policy* for
> multi-agent work. Read it that way, not as a description of a CLI feature.

Rules governing multi-agent teamwork previews (`teamwork-preview`), subagent delegation, and artifact lifecycle management across the AGY harness control plane.

## 1. Multi-Agent Teamwork Previews
- **Agent Scope & Write Boundaries**: Every delegated subagent MUST operate strictly within its assigned folder inside `.agents/` or designated file scope. Overlapping writes outside owned boundaries are prohibited.
- **Role Activation**: Subagent tasks MUST explicitly specify role assignments (e.g. `implementer`, `qa`, `specialist`). Primary execution follows the `implementer` role unless specified.
- **Parallel Dispatch**: Dispatch tasks in parallel to independent subagent instances whenever possible to minimize latency.

## 2. Subagent Delegation Workflow
1. **Prompt Crafting**: Formulate structured prompts specifying target working directory, dispatch path, ownership boundary, objectives, and acceptance criteria.
2. **Draft Recording**: Record prompt parameters and task specifications in the subagent's `DISPATCH.md`.
3. **Objective Verification**: Define clear, objective, and independently executable verification commands (`make test`, `make check`, self-tests).
4. **Integrity Mode**: Maintain real state and logic. Hardcoded outputs, dummy implementations, or fabricated verification logs are strictly forbidden.

## 3. Artifact Lifecycle Management
- **Directory Structure**: All agent metadata (plans, progress heartbeats, briefing files, handoffs) MUST be placed in `.agents/<agent_name>/`. Source code, tests, or runtime data must NEVER be placed in `.agents/`.
- **BRIEFING.md**: Every subagent MUST maintain `BRIEFING.md` as its working memory, preserving append-only identity and key constraint sections.
- **Progress Heartbeat**: `progress.md` MUST be updated after every meaningful step to serve as the liveness heartbeat.
- **Handoff Report**: Every completed task MUST produce a 5-component `handoff.md` (Observation, Logic Chain, Caveats, Conclusion, Verification Method).

## 4. Verification & Clean Handover
- Execute full build, test, lint, and convention checks before declaring task completion.
- Send handoff summary and `handoff.md` path to parent agent via `send_message`.

## 5. Execution Runtime & Anti-Polling Guards
- **Anti-Polling Invariant**: NEVER poll or loop on `manage_task(Action='status')` or any task monitoring loop. Antigravity's runtime is reactive and event-driven; the system automatically delivers a `<SYSTEM_MESSAGE>` when background tasks or subagents finish. Always conclude your turn or proceed to disjoint work immediately after launching.
- **Canonical Workspace Grounding**: Always verify the active working directory and branch with `git branch --show-current` and `git status --short --branch`. If executing in a managed worktree, ensure changes landed on `origin/main` are synchronized back to canonical sibling repositories (`~/hairglasses-studio/<repo>`).
- **Non-Interactive Shell Execution**: Always execute shell commands with non-interactive flags and pipe stdin from `/dev/null` (e.g., `cmd < /dev/null`, `PAGER=cat`, `git --no-pager`) to avoid stalling on interactive terminal prompts.
- **Database & Test Isolation**: Unit test suites must never depend on remote WAN databases (e.g. Neon) to prevent quota exhaustion (`SQLSTATE 53000`) and network latency. Test fixtures must use local ephemeral SQLite / test-scoped mock DB / local PG with trust authentication.
