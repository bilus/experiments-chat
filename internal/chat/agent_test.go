package chat

import "testing"

func TestStaticProviderReturnsOnlyResponseText(t *testing.T) {
	provider := StaticProvider{Response: "I don't understand."}
	conversation := Conversation{
		{Role: RoleUser, Text: "hello"},
		{Role: RoleAssistant, Text: "previous response"},
	}

	response, err := provider.Ask("system prompt", conversation, "new question")
	if err != nil {
		t.Fatalf("ask static provider: %v", err)
	}

	if response != "I don't understand." {
		t.Fatalf("expected response text only, got %q", response)
	}
}

func TestAskAppendsUserInputAndProviderResponseToConversation(t *testing.T) {
	agent := NewAgent("system prompt", StaticProvider{Response: "I don't understand."})

	agent, err := Ask(agent, "hello")
	if err != nil {
		t.Fatalf("ask agent: %v", err)
	}

	if len(agent.Conversation) != 2 {
		t.Fatalf("expected 2 conversation messages, got %d", len(agent.Conversation))
	}

	if agent.Conversation[0] != (Message{Role: RoleUser, Text: "hello"}) {
		t.Fatalf("expected user message, got %#v", agent.Conversation[0])
	}

	if agent.Conversation[1] != (Message{Role: RoleAssistant, Text: "I don't understand."}) {
		t.Fatalf("expected assistant message, got %#v", agent.Conversation[1])
	}
}
