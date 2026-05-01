package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bilus/experiments-chat/internal/chat"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
)

type terminalChatFeature struct {
	model     Model
	userInput string
}

func (f *terminalChatFeature) theTerminalChatApplicationIsRunning() error {
	provider := chat.StaticProvider{Response: "I don't understand."}
	f.model = NewModel(chat.NewAgent("", provider))
	return nil
}

func (f *terminalChatFeature) theUserEnters(text string) error {
	f.userInput = text

	updated, _ := f.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)})
	f.model = updated.(Model)

	updated, _ = f.model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f.model = updated.(Model)

	return nil
}

func (f *terminalChatFeature) theTerminalChatShouldRender(expected string) error {
	view := f.model.View()
	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected terminal chat to render %q, got %q", expected, view)
	}

	conversation := f.model.agent.Conversation
	if len(conversation) != 2 {
		return fmt.Errorf("expected agent conversation to contain 2 messages, got %d", len(conversation))
	}

	if conversation[0].Role != chat.RoleUser || conversation[0].Text != f.userInput {
		return fmt.Errorf("expected first message to be user input %q, got role %q text %q", f.userInput, conversation[0].Role, conversation[0].Text)
	}

	if conversation[1].Role != chat.RoleAssistant || conversation[1].Text != expected {
		return fmt.Errorf("expected second message to be assistant response %q, got role %q text %q", expected, conversation[1].Role, conversation[1].Text)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			feature := &terminalChatFeature{}

			sc.Step(`^the terminal chat application is running$`, feature.theTerminalChatApplicationIsRunning)
			sc.Step(`^the user enters "([^"]*)"$`, feature.theUserEnters)
			sc.Step(`^the terminal chat should render "([^"]*)"$`, feature.theTerminalChatShouldRender)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("godog feature suite failed")
	}
}
