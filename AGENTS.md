# AGENTS.md

This file is the repo-local briefing for agents contributing to `experiments-chat`.

Read the upstream environment guide first:

https://wiki.bilus.dev/agents

That page owns the global style and patch-stack rules. This file points to the detailed project-local guides an agent should read before working in a specific area.

## Non-Negotiable Rules

- Do not write em dashes in chat replies, docs, issues, commit messages, code comments, or PR text.
- Work on the `development` branch unless the user explicitly says otherwise.
- Use `gps` and the patch-stack workflow for PR management.
- Do not create separate feature branches for normal work in this repo.
- Keep each patch independent: 1 logical change = 1 commit = 1 PR.
- Amend the active patch instead of piling fixup commits on top.
- Never commit without explicit user request.
- Do not touch unrelated untracked files. `.DS_Store` and local artifacts may exist.
- Use devbox and direnv for project commands.

## Read Before Working

- Requirements or GitHub issues: read [docs/agent-requirements.md](docs/agent-requirements.md).
- TDD, Gherkin, Godog, or Go tests: read [docs/agent-testing.md](docs/agent-testing.md).
- Git, `gps`, commits, or PR flow: read [docs/agent-workflow.md](docs/agent-workflow.md).
- Devbox, direnv, or local commands: read [docs/agent-environment.md](docs/agent-environment.md).
- Agent and provider boundaries: read [docs/agent-provider-contract.md](docs/agent-provider-contract.md).
- ACP spike context: read [docs/compatibility-findings.md](docs/compatibility-findings.md) and [docs/provider-launch.md](docs/provider-launch.md).
- Logging expectations: read [docs/logging-and-observability.md](docs/logging-and-observability.md).
- Stage boundaries: read [docs/planning-guardrails.md](docs/planning-guardrails.md).

## Default Working Loop

1. Read the relevant detailed guide above.
2. Inspect the current issue, project status, code, and tests.
3. Write or update the requirement-level Gherkin test first when behavior changes.
4. Run `direnv exec . go test ./...` and confirm the expected red state.
5. Implement the smallest useful change.
6. Run `direnv exec . go test ./...` again.
7. Stage only intended files.
8. Amend or commit only when the user has explicitly asked.

## Documentation

Keep docs small and focused.

Use `docs/` for project-local notes that guide requirements and implementation. When docs summarize ACP spike findings, prefer linking or naming the source artifact rather than copying large reports wholesale.

