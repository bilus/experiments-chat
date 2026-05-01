package tui

import (
	"bytes"
	"testing"

	"github.com/bilus/experiments-chat/internal/chat"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTerminalChatRendersPrompt(t *testing.T) {
	provider := chat.StaticProvider{Response: "I don't understand."}
	agent := chat.NewAgent("", provider)
	program := tea.NewProgram(NewModel(agent), tea.WithoutRenderer(), tea.WithInput(bytes.NewReader(nil)))

	model, err := program.Run()
	if err != nil {
		t.Fatalf("run program: %v", err)
	}

	view := model.(Model).View()
	if view != "You: " {
		t.Fatalf("expected minimal prompt, got %q", view)
	}
}
