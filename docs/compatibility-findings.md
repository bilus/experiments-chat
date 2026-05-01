# ACP Compatibility Findings Relevant To Requirements

The most useful comprehensive baseline run from the spike was the full provider run on 2026-04-23.

It used the tiny Go workspace scenario: inspect the workspace, add `Mul(a, b int) int`, add tests, run `go test ./...`, ask a follow-up question, and attempt session load or resume.

Source artifacts:

- baseline: `llm-wiki/tmp/runs/20260423-220322.967291000-baseline`
- delegated worker: `llm-wiki/tmp/runs/20260423-222624.707952000-delegated`
- Codex permission probe: `llm-wiki/tmp/runs/20260424-093541.659677000-permission-probe`

## Baseline Matrix Summary

| Feature | Claude | Copilot | Codex | Requirement impact |
| --- | --- | --- | --- | --- |
| Capability discovery | Works | Works | Works | Keep capability discovery in the provider adapter. Do not hard-code all provider behavior. |
| Session creation | Works | Works | Works | Terminal chat can start from a session abstraction. |
| Initial prompt | Works | Works | Works | The minimal chat path is a good first requirement. |
| Streaming output | Works | Works | Works | Later UI work should expose streaming, even if iteration 1 renders final text only. |
| Multi-turn follow-up | Works | Works | Works | Conversation continuity is viable, but app-level logs should still record messages. |
| Session resume/load | Works | Works with adaptation | Works | Resume must be modeled as provider-sensitive. Copilot reported an already-loaded session as an error-shaped response. |
| Approval and permission flow | Works | Works | Not verified | Permission behavior must stay explicit in requirements. Codex needed a separate probe. |
| Workspace file edits | Not verified in that row for Claude | Works | Works | File edits are viable, but row-level evidence must be preserved instead of inferred from final status. |
| Command execution | Works with adaptation | Works with adaptation | Works | Command execution needs provider-specific approval and failure handling. |
| JSONL event log | Works | Works | Works | Logging is not optional for production use. |
| Human transcript rendering | Works | Works | Works | Human-readable run rendering should be built from JSONL logs. |
| Scenario outcome | Works | Works with adaptation | Works with adaptation | ACP is promising, but the app needs an adapter layer rather than direct provider branching in UI code. |

## Important Provider Deviations

- Claude reached the baseline through the ACP wrapper over the local Claude CLI.
- Copilot resume can report that a session is already loaded. Treat this as an adaptation case, not an automatic hard failure.
- Copilot command execution may be blocked by approval. The UI and logs should surface this clearly.
- Codex completed baseline edits and tests, but in one run the follow-up repeated the first response. Requirements should not assume perfect answer quality from ACP compatibility alone.
- Codex permission flow was not verified by the baseline because no approval prompts surfaced in that run.

## Permission Probe

A later Codex-only permission probe observed ACP approval, but the scenario was still partial because the normal in-workspace task did not complete afterward.

Requirement impact:

- treat permission flow as independently testable from normal chat
- check disk state and logs, not only model text
- do not mark all Codex permissions solved from one approval event

## Delegated-Worker Scenario

The delegated-worker scenario is closer to the planned supervisor shape, where this app asks a provider to act as a worker and report back.

It should stay separate from baseline ACP compatibility because it also measures prompt following, reporting discipline, and model behavior.

Requirement impact:

- use a separate test scenario for worker-style prompts
- keep baseline chat requirements simple
- do not mix delegated-worker scores into core ACP transport compatibility
