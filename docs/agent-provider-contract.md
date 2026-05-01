# Agent And Provider Contract

The prototype should keep a clear boundary between app state and ACP transport details.

## Current Prototype Contract

The current minimal contract is:

```go
type Provider interface {
	Ask(systemPrompt string, conversation Conversation, text string) (string, error)
}
```

This is an application-facing contract, not a direct ACP protocol type.

The key rule is:

- `Provider.Ask` returns only the assistant response text.
- `Agent.Ask` owns appending the user message and assistant response to the application conversation.

This matches the direction we want for ACP compatibility: provider adapters may use sessions, streaming events, and callbacks internally, but the first chat-facing contract should not return a whole rewritten conversation.

## Why This Matters

ACP providers can keep their own session state. The app still needs its own conversation record for:

- terminal rendering
- run logs
- replay
- future transcript export
- tests that do not depend on a provider implementation

If provider adapters returned full conversations, every provider quirk would leak upward. Keeping the return value to response text makes the Agent responsible for consistent app-level state.

## Near-Term Implications

The `StaticProvider` is a test and bootstrapping provider. It should stay small.

The first real ACP provider can implement the same `Provider` interface by:

- creating or loading an ACP session
- sending the prompt
- collecting final assistant text from streamed events
- returning that final response text
- leaving conversation mutation to `Agent.Ask`

When streaming becomes a requirement, this interface should evolve deliberately. A likely next shape is an event-based provider API rather than stretching `Ask` too far.

