# Agent Workflow Guide

Use this guide when staging, committing, amending, or preparing PRs.

## Upstream Reference

The authoritative global workflow is the OtterWiki AGENTS page:

https://wiki.bilus.dev/agents

This repo follows the same GPS patch-stack discipline.

## Branch Model

Work on `development`.

Do not create feature branches for normal issue work. `gps` creates review branches when patches are prepared for review.

## Patch Model

Rules:

- 1 logical change = 1 commit = 1 PR.
- Each commit must be independently understandable.
- Each commit should be independently testable when possible.
- Avoid editing the same file from multiple stacked commits.
- If two patches need the same file, fold or reorder the work before review.
- Amend the active patch instead of creating fixup commits.

## Commit Rules

Never commit without explicit user request.

When the user has asked for a commit or amend:

```bash
direnv exec . go test ./...
git status --short
git add <intended files>
git commit --amend --no-edit
```

Use a normal `git commit` only when creating a new independent patch.

## GPS Rules

Use `gps` for patch-stack and PR management.

Follow the upstream convention from OtterWiki AGENTS:

- use explicit review branch names
- keep branch names short and descriptive
- rely on CI for catch-up when local untracked files make isolation checks noisy

Do not use `gps` to hide messy local state. Clean the intended patch first.

## Local State

The repo may contain unrelated untracked files such as `.DS_Store`.

Do not stage, delete, or rewrite unrelated files unless the user asks.

