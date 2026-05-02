package test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/coder/acp-go-sdk"
)

const helperProcessEnvName = "FAKE_AGENT"

func UseFakeAgent(t *testing.T, commandOverrideEnvName string) {
	t.Helper()

	command := []string{os.Args[0], "-test.run=TestFakeAgent", "--"}

	commandJSON, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("marshal helper command: %v", err)
	}

	t.Setenv(commandOverrideEnvName, string(commandJSON))
	t.Setenv(helperProcessEnvName, "1")
}

func RunFakeAgent(t *testing.T) {
	t.Helper()

	if os.Getenv(helperProcessEnvName) != "1" {
		return
	}

	agent := &fakeAgent{}
	connection := acp.NewAgentSideConnection(agent, os.Stdout, os.Stdin)
	agent.connection = connection
	connection.SetLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	<-connection.Done()
	os.Exit(0)
}

type fakeAgent struct {
	connection *acp.AgentSideConnection
}

func (*fakeAgent) Authenticate(context.Context, acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	return acp.AuthenticateResponse{}, nil
}

func (*fakeAgent) Initialize(context.Context, acp.InitializeRequest) (acp.InitializeResponse, error) {
	return acp.InitializeResponse{
		ProtocolVersion: acp.ProtocolVersionNumber,
		AgentCapabilities: acp.AgentCapabilities{
			LoadSession: false,
		},
	}, nil
}

func (*fakeAgent) Cancel(context.Context, acp.CancelNotification) error {
	return nil
}

func (*fakeAgent) ListSessions(context.Context, acp.ListSessionsRequest) (acp.ListSessionsResponse, error) {
	return acp.ListSessionsResponse{}, errors.New("not implemented")
}

func (*fakeAgent) NewSession(context.Context, acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	return acp.NewSessionResponse{SessionId: "helper-session"}, nil
}

func (a *fakeAgent) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	if len(params.Prompt) != 1 || params.Prompt[0].Text == nil {
		return acp.PromptResponse{}, errors.New("expected one text prompt")
	}

	text := params.Prompt[0].Text.Text
	for _, chunk := range []string{"helper heard: ", text} {
		if err := a.connection.SessionUpdate(ctx, acp.SessionNotification{
			SessionId: params.SessionId,
			Update:    acp.UpdateAgentMessageText(chunk),
		}); err != nil {
			return acp.PromptResponse{}, err
		}
	}

	return acp.PromptResponse{StopReason: acp.StopReasonEndTurn}, nil
}

func (*fakeAgent) SetSessionConfigOption(context.Context, acp.SetSessionConfigOptionRequest) (acp.SetSessionConfigOptionResponse, error) {
	return acp.SetSessionConfigOptionResponse{}, errors.New("not implemented")
}

func (*fakeAgent) SetSessionMode(context.Context, acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	return acp.SetSessionModeResponse{}, errors.New("not implemented")
}
