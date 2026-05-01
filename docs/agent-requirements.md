# Agent Requirements Guide

Use this guide when creating, reviewing, or implementing requirements.

## Storage Model

Requirements are captured as GitHub issues.

Rules:

- Each requirement gets its own issue.
- The issue number is the requirement ID.
- Requirement text uses EARS syntax.
- Keep issue bodies small and focused.
- Put requirement text under `# Requirement`.
- Put acceptance tests under `# Test Plan` using Gherkin.
- Use GraphQL for issue creation and GitHub Project updates when doing bulk requirement work.
- Track requirements in the GitHub Project at `https://github.com/users/bilus/projects/3/views/1`.

Project columns:

- Add every requirement issue to `Requirements`.
- Move only the current iteration's work to `Todo`.
- Move the actively worked issue to `In Progress`.
- Move issues to `Done` only after implementation and tests are complete.

## EARS Format

EARS stands for Easy Approach to Requirements Syntax.

Prefer one of these forms:

| Pattern | Form | Use when |
| --- | --- | --- |
| Ubiquitous | The system shall `<response>`. | The behavior is always true. |
| Event-driven | When `<trigger>`, the system shall `<response>`. | A behavior happens after an event. |
| State-driven | While `<state>`, the system shall `<response>`. | A behavior applies only in a state. |
| Unwanted behavior | If `<unwanted condition>`, then the system shall `<response>`. | The system handles an error or exceptional case. |
| Optional feature | Where `<feature is included>`, the system shall `<response>`. | A behavior applies only when a feature exists. |
| Complex | While `<state>`, when `<trigger>`, the system shall `<response>`. | Both state and event matter. |

Good requirements are:

- observable
- testable
- singular
- provider-neutral unless the requirement is explicitly provider-specific
- written from the system's point of view

Avoid:

- implementation details unless the requirement is explicitly technical
- vague verbs such as support, handle, improve, or manage without a concrete observable result
- multiple behaviors joined by and
- hidden acceptance criteria not reflected in the Gherkin test plan

## Examples

Ubiquitous:

```text
The system shall provide a terminal chat application runnable from a single binary.
```

Event-driven:

```text
When the user submits a chat message, the system shall send the message to the configured agent provider.
```

State-driven:

```text
While an agent request is in progress, the system shall render that the assistant is working.
```

Unwanted behavior:

```text
If the configured provider returns an error, then the system shall render the error without losing the user's conversation history.
```

Optional feature:

```text
Where JSONL logging is enabled, the system shall write each provider interaction as a structured log event.
```

Complex:

```text
While a chat session is active, when the provider requests permission, the system shall render the permission request and wait for a user decision.
```

## Issue Body Template

Use this shape for simple requirement issues:

````markdown
# Requirement
The system shall provide a terminal chat application runnable from a single binary.

# Test Plan
```gherkin
Feature: Minimal terminal chat prototype
  Scenario: Respond to any user input with the PoC response
    Given the terminal chat application is running
    When the user enters "hello"
    Then the terminal chat should render "I don't understand."
```
````

Keep the issue body to the requirement and test plan unless the user asks for more.

## Implementation Linkage

When implementing a requirement:

- read the issue first
- keep the issue number visible in your notes and final response
- add or update the matching `.feature` file before production code
- keep Gherkin wording close to the issue test plan
- add unit tests for lower-level behavior that the Gherkin scenario depends on
- do not move the issue to `Done` until `direnv exec . go test ./...` passes

