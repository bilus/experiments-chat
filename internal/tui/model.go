package tui

import (
	"strings"

	"github.com/bilus/experiments-chat/internal/chat"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	agent chat.Agent
	input string
	err   error
}

func NewModel(agent chat.Agent) Model {
	return Model{agent: agent}
}

func (Model) Init() tea.Cmd {
	return tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.Type {
	case tea.KeyRunes:
		m.input += string(key.Runes)
	case tea.KeyEnter:
		if m.input != "" {
			agent, err := chat.Ask(m.agent, m.input)
			if err != nil {
				m.err = err
				return m, nil
			}

			m.agent = agent
			m.input = ""
		}
	}

	return m, nil
}

func (m Model) View() string {
	var lines []string
	for _, message := range m.agent.Conversation {
		if message.Role == chat.RoleAssistant {
			lines = append(lines, message.Text)
		}
	}

	if m.err != nil {
		lines = append(lines, m.err.Error())
	}

	lines = append(lines, "You: "+m.input)
	return strings.Join(lines, "\n")
}
