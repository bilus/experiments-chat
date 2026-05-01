# Agent Testing Guide

Use this guide when adding or changing behavior.

## Test-Driven Development

Use TDD for every behavior change.

The flow:

1. Add or update the Gherkin test plan in the requirement issue.
2. Add the matching `.feature` file to the repo.
3. Implement Godog step definitions with `github.com/cucumber/godog`.
4. Run `direnv exec . go test ./...` and confirm the test fails for the expected reason.
5. Add the smallest production code that satisfies the failing test.
6. Add focused Go unit tests for domain behavior, edge cases, and provider contracts.
7. Run `direnv exec . go test ./...` again and confirm it passes.
8. Amend the active patch when the user has requested a commit.

Do not skip the red step. A test that never failed did not prove the requirement.

## Godog

Godog is for requirement-level Gherkin acceptance tests.

Use it when a requirement issue has a `# Test Plan` section with Gherkin.

Recommended layout:

```text
internal/<area>/features/<requirement-name>.feature
internal/<area>/<requirement-name>_feature_test.go
```

The `.feature` file should mirror the issue test plan unless the issue was intentionally refined.

## Go Unit Tests

Go unit tests are for smaller behavior such as:

- agent state transitions
- provider contract guarantees
- error handling
- rendering decisions
- adapter normalization

Use unit tests to pin down details that would make a Gherkin scenario noisy or brittle.

## Verification Command

Use:

```bash
direnv exec . go test ./...
```

Run it after the red test is written, after the implementation, and before staging or amending.

Do not substitute `go build` unless the user asks for it.

## Formatting

Use `gofumpt`, not `gofmt`, for Go formatting.

```bash
direnv exec . gofumpt -w <go files>
```
