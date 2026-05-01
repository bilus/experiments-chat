package chat

type StaticProvider struct {
	Response string
}

func (p StaticProvider) Ask(systemPrompt string, conversation Conversation, text string) (string, error) {
	return p.Response, nil
}
