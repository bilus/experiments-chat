package tui

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bilus/experiments-chat/internal/chat"
	"github.com/bilus/experiments-chat/internal/claude"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	quitShortcutHelp       = "Press Ctrl+D to quit."
	claudeReadyMessage     = "Claude is ready."
	claudeConnectingStatus = "Connecting to Claude..."
	claudeNotReadyMessage  = "Claude is still connecting."
	claudeStartupTimeout   = 30 * time.Second
)

type Model struct {
	agent            chat.Agent
	prompt           textinput.Model
	err              error
	status           string
	initCmd          tea.Cmd
	startupCancel    context.CancelFunc
	providerShutdown func() error
	providerReady    bool
	asking           bool
}

type claudeReadyMsg struct {
	connection *claude.Connection
}

type claudeStartupFailedMsg struct {
	err error
}

type claudeStartupCanceledMsg struct{}

type startClaudeMsg struct{}

type agentAskedMsg struct {
	agent chat.Agent
	err   error
}

func NewModel(agent chat.Agent) Model {
	prompt := textinput.New()
	prompt.Prompt = "You: "
	prompt.Focus()

	return Model{
		agent:         agent,
		prompt:        prompt,
		providerReady: true,
	}
}

func NewWithClaude() Model {
	model := NewModel(chat.NewAgent("", nil))
	model.status = claudeConnectingStatus
	model.initCmd = startClaude
	model.providerReady = false
	return model
}

func (m Model) Init() tea.Cmd {
	return m.initCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startClaudeMsg:
		return m.startClaude()
	case claudeReadyMsg:
		m.startupCancel = nil
		m.status = claudeReadyMessage
		m.agent = chat.WithProvider(m.agent, msg.connection)
		m.providerShutdown = msg.connection.Close
		m.providerReady = true
		return m, nil
	case claudeStartupFailedMsg:
		m.startupCancel = nil
		m.status = ""
		m.err = msg.err
		return m, nil
	case claudeStartupCanceledMsg:
		m.startupCancel = nil
		return m, nil
	case agentAskedMsg:
		m.asking = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.agent = msg.agent
		return m, nil
	}

	key, ok := msg.(tea.KeyMsg)
	if ok {
		switch key.Type {
		case tea.KeyCtrlD:
			m.cancelStartup()
			return m, tea.Quit
		case tea.KeyEnter:
			text := m.prompt.Value()
			if text != "" {
				if !m.providerReady {
					m.err = errors.New(claudeNotReadyMessage)
					return m, nil
				}
				if m.asking {
					return m, nil
				}

				m.prompt.Reset()
				m.asking = true
				return m, askAgent(m.agent, text)
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
	if m.status != "" {
		lines = append(lines, m.status)
	}

	for _, message := range m.agent.Conversation {
		lines = append(lines, renderConversationMessage(message))
	}

	if m.err != nil {
		lines = append(lines, m.err.Error())
	}

	lines = append(lines, m.prompt.View())
	return strings.Join(lines, "\n")
}

func (m *Model) Close() error {
	return m.stopProvider()
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

func askAgent(agent chat.Agent, text string) tea.Cmd {
	return func() tea.Msg {
		agent, err := chat.Ask(agent, text)
		return agentAskedMsg{
			agent: agent,
			err:   err,
		}
	}
}

func startClaude() tea.Msg {
	return startClaudeMsg{}
}

func (m Model) startClaude() (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), claudeStartupTimeout)
	m.startupCancel = cancel

	return m, func() tea.Msg {
		connection, err := claude.Initialize(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return claudeStartupCanceledMsg{}
			}
			return claudeStartupFailedMsg{err: err}
		}

		return claudeReadyMsg{connection: connection}
	}
}

func (m *Model) cancelStartup() {
	if m.startupCancel == nil {
		return
	}

	m.startupCancel()
	m.startupCancel = nil
}

func (m *Model) stopProvider() error {
	if m.providerShutdown == nil {
		return nil
	}

	err := m.providerShutdown()
	m.err = err
	m.providerShutdown = nil
	return err
}
