package chat

import "errors"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role Role
	Text string
}

type Conversation []Message

type Provider interface {
	Ask(systemPrompt string, conversation Conversation, text string) (string, error)
}

type Agent struct {
	SystemPrompt string
	Conversation Conversation
	provider     Provider
}

func NewAgent(systemPrompt string, provider Provider) Agent {
	return Agent{
		SystemPrompt: systemPrompt,
		provider:     provider,
	}
}

func WithProvider(agent Agent, provider Provider) Agent {
	agent.provider = provider
	return agent
}

func Ask(agent Agent, text string) (Agent, error) {
	if agent.provider == nil {
		return agent, errors.New("chat provider is required")
	}

	response, err := agent.provider.Ask(agent.SystemPrompt, agent.Conversation, text)
	if err != nil {
		return agent, err
	}

	agent.Conversation = append(
		agent.Conversation,
		Message{Role: RoleUser, Text: text},
		Message{Role: RoleAssistant, Text: response},
	)

	return agent, nil
}
