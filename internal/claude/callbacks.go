package claude

import (
	"context"
	"errors"
	"os/exec"

	acp "github.com/coder/acp-go-sdk"
)

type callbacks struct {
	streamingResponse *streamingResponse
}

func (callbacks) ReadTextFile(context.Context, acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	return acp.ReadTextFileResponse{}, errors.New("ACP file reads are not supported yet")
}

func (callbacks) WriteTextFile(context.Context, acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	return acp.WriteTextFileResponse{}, errors.New("ACP file writes are not supported yet")
}

func (callbacks) RequestPermission(context.Context, acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	return acp.RequestPermissionResponse{
		Outcome: acp.RequestPermissionOutcome{
			Cancelled: &acp.RequestPermissionOutcomeCancelled{},
		},
	}, nil
}

func (c callbacks) SessionUpdate(ctx context.Context, notification acp.SessionNotification) error {
	if c.streamingResponse == nil {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return ctx.Err()
	}
	c.streamingResponse.Append(notification)
	return nil
}

func (callbacks) CreateTerminal(context.Context, acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	return acp.CreateTerminalResponse{}, errors.New("ACP terminal creation is not supported yet")
}

func (callbacks) KillTerminal(context.Context, acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	return acp.KillTerminalResponse{}, errors.New("ACP terminal kill is not supported yet")
}

func (callbacks) TerminalOutput(context.Context, acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	return acp.TerminalOutputResponse{}, errors.New("ACP terminal output is not supported yet")
}

func (callbacks) ReleaseTerminal(context.Context, acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	return acp.ReleaseTerminalResponse{}, errors.New("ACP terminal release is not supported yet")
}

func (callbacks) WaitForTerminalExit(context.Context, acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	return acp.WaitForTerminalExitResponse{}, errors.New("ACP terminal wait is not supported yet")
}

func normalizeProcessWaitError(err error) error {
	if err == nil {
		return nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ProcessState != nil && exitErr.ProcessState.ExitCode() == 0 {
		return nil
	}

	return err
}

func shorten(text string, limit int) string {
	if len(text) <= limit {
		return text
	}

	return text[:limit] + "..."
}
