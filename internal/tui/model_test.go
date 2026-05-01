package tui

import (
	"strings"
	"testing"

	"github.com/bilus/experiments-chat/internal/chat"
)

func TestTerminalChatRendersPrompt(t *testing.T) {
	provider := chat.StaticProvider{Response: "I don't understand."}
	agent := chat.NewAgent("", provider)

	view := NewModel(agent).View()
	for _, expected := range []string{"Press Ctrl+D to quit.", "You: "} {
		if !strings.Contains(view, expected) {
			t.Fatalf("expected prompt view to contain %q, got %q", expected, view)
		}
	}
}
