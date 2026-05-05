package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/bilus/experiments-chat/internal/chat"
	tea "github.com/charmbracelet/bubbletea"
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

func TestSubmittingPromptReturnsCommandWithoutBlockingOnProvider(t *testing.T) {
	releaseProvider := make(chan struct{})
	defer close(releaseProvider)

	provider := blockingProvider{release: releaseProvider}
	model := NewModel(chat.NewAgent("", provider))

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
	model = updated.(Model)

	done := make(chan tea.Cmd, 1)
	go func() {
		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		done <- cmd
	}()

	select {
	case cmd := <-done:
		if cmd == nil {
			t.Fatal("expected prompt submission to return an ask command")
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("prompt submission blocked on provider")
	}
}

func TestSubmittingPromptWhileProviderIsBusyDoesNotStartStaleRequest(t *testing.T) {
	releaseProvider := make(chan struct{})
	defer close(releaseProvider)

	provider := blockingProvider{release: releaseProvider}
	model := NewModel(chat.NewAgent("", provider))

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("first")})
	model = updated.(Model)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("expected first prompt submission to start an ask command")
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("second")})
	model = updated.(Model)
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("expected second prompt submission to be ignored while provider is busy")
	}
}

func TestSubmittingPromptWhileClaudeConnectsShowsConnectionMessage(t *testing.T) {
	model := NewWithClaude()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
	model = updated.(Model)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	if cmd != nil {
		t.Fatal("expected prompt submission during Claude startup not to start an ask command")
	}
	if !strings.Contains(model.View(), claudeNotReadyMessage) {
		t.Fatalf("expected connection message %q, got %q", claudeNotReadyMessage, model.View())
	}
}

func TestModelCloseStopsProvider(t *testing.T) {
	model := NewModel(chat.NewAgent("", chat.StaticProvider{Response: "unused"}))
	closed := false
	model.providerShutdown = func() error {
		closed = true
		return nil
	}

	model.Close()

	if !closed {
		t.Fatal("expected model close to stop provider")
	}
	if model.providerShutdown != nil {
		t.Fatal("expected model close to clear provider shutdown")
	}
}

type blockingProvider struct {
	release <-chan struct{}
}

func (p blockingProvider) Ask(string, chat.Conversation, string) (string, error) {
	<-p.release
	return "done", nil
}
