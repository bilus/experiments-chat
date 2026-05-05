package claude

import (
	"context"
	"testing"
	"time"

	"github.com/bilus/experiments-chat/internal/test"
)

const helperProcessEnvName = "EXPERIMENTS_CHAT_CLAUDE_ACP_HELPER"

func TestFakeAgent(t *testing.T) {
	test.RunFakeAgent(t)
}

func TestInitializeConnectsToConfiguredACPWrapper(t *testing.T) {
	test.UseFakeAgent(t, commandOverrideEnvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	connection, err := Initialize(ctx)
	if err != nil {
		t.Fatalf("initialize Claude ACP wrapper: %v", err)
	}
	defer func() {
		if err := connection.Close(); err != nil {
			t.Fatalf("close Claude ACP wrapper: %v", err)
		}
	}()
}

func TestConnectionAskReturnsAgentMessageText(t *testing.T) {
	test.UseFakeAgent(t, commandOverrideEnvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	connection, err := Initialize(ctx)
	if err != nil {
		t.Fatalf("initialize Claude ACP wrapper: %v", err)
	}
	defer func() {
		if err := connection.Close(); err != nil {
			t.Fatalf("close Claude ACP wrapper: %v", err)
		}
	}()

	response, err := connection.Ask("", nil, "hello")
	if err != nil {
		t.Fatalf("ask Claude ACP wrapper: %v", err)
	}
	if response != "helper heard: hello" {
		t.Fatalf("expected provider to return assistant response text, got %q", response)
	}
}
