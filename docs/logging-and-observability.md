# Logging And Observability

The ACP spike treated logging as a first-class deliverable. The chat prototype should keep that direction.

## Source Of Truth

The machine-readable log should be newline-delimited JSON.

Human-readable transcripts should be rendered from that JSONL log, not maintained as a separate source of truth.

## Events To Capture

Real ACP provider runs should eventually log:

- run ID and provider name
- provider command and version when available
- session creation, load, and resume attempts
- prompts submitted by the user or system
- streamed assistant events
- tool or file callbacks
- approval requests and selected outcomes
- command execution events when visible
- errors and provider stderr summaries
- cancellation and retry attempts
- final verification results when a scenario includes verification

## Requirement Impact

For the minimal terminal chat requirement, logs can be deferred if the issue is intentionally scoped to rendering a single response.

For production-readiness requirements, logs should be mandatory before real Claude, Copilot, or Codex providers are considered usable.

The UI should also have a simple way to render logs into a readable transcript so a user can inspect what happened without reading raw JSONL.

