package tui

import (
	"strings"

	"github.com/bilus/experiments-chat/internal/chat"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const quitShortcutHelp = "Press Ctrl+D to quit."

type Model struct {
	agent  chat.Agent
	prompt textinput.Model
	err    error
}

func NewModel(agent chat.Agent) Model {
	prompt := textinput.New()
	prompt.Prompt = "You: "
	prompt.Focus()

	return Model{
		agent:  agent,
		prompt: prompt,
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if ok {
		switch key.Type {
		case tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyEnter:
			text := m.prompt.Value()
			if text != "" {
				agent, err := chat.Ask(m.agent, text)
				if err != nil {
					m.err = err
					return m, nil
				}

				m.agent = agent
				m.prompt.Reset()
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	lines := []string{quitShortcutHelp}
	for _, message := range m.agent.Conversation {
		lines = append(lines, renderConversationMessage(message))
	}

	if m.err != nil {
		lines = append(lines, m.err.Error())
	}

	lines = append(lines, m.prompt.View())
	return strings.Join(lines, "\n")
}

func renderConversationMessage(message chat.Message) string {
	switch message.Role {
	case chat.RoleUser:
		return "You: " + message.Text
	case chat.RoleAssistant:
		return "Assistant: " + message.Text
	default:
		return string(message.Role) + ": " + message.Text
	}
}
