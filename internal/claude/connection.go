package claude

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bilus/experiments-chat/internal/buf"
	"github.com/bilus/experiments-chat/internal/chat"
	acp "github.com/coder/acp-go-sdk"
)

const (
	commandOverrideEnvName = "EXPERIMENTS_CHAT_CLAUDE_ACP_COMMAND"
	promptTimeout          = 5 * time.Minute
	shutdownGracePeriod    = 2 * time.Second
	shutdownKillWait       = 2 * time.Second
)

type Connection struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	conn      *acp.ClientSideConnection
	stderr    *buf.Buffer
	waitCh    chan error
	sessionID acp.SessionId
	response  *streamingResponse
	askMu     sync.Mutex
	closeOnce sync.Once
}

func Initialize(ctx context.Context) (*Connection, error) {
	connection, err := start()
	if err != nil {
		return nil, err
	}

	if err := connection.initialize(ctx); err != nil {
		closeErr := connection.Close()
		if closeErr != nil {
			err = fmt.Errorf("%w; close Claude ACP wrapper: %v", err, closeErr)
		}
		return nil, err
	}

	return connection, nil
}

func start() (*Connection, error) {
	command, err := command()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(command[0], command[1:]...)
	stderr := &buf.Buffer{}
	cmd.Stderr = stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open Claude ACP stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("open Claude ACP stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start Claude ACP wrapper %q: %w", strings.Join(command, " "), err)
	}

	response := &streamingResponse{}
	conn := acp.NewClientSideConnection(callbacks{streamingResponse: response}, stdin, stdout)
	conn.SetLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
		close(waitCh)
	}()

	return &Connection{
		cmd:      cmd,
		stdin:    stdin,
		conn:     conn,
		stderr:   stderr,
		waitCh:   waitCh,
		response: response,
	}, nil
}

func (c *Connection) initialize(ctx context.Context) error {
	_, err := c.conn.Initialize(ctx, acp.InitializeRequest{
		ProtocolVersion: acp.ProtocolVersionNumber,
		ClientInfo: &acp.Implementation{
			// TODO(bilus): Pass app name and version.
			Name:    "experiments-chat",
			Version: "dev",
		},
		ClientCapabilities: acp.ClientCapabilities{},
	})
	if err != nil {
		return c.wrapProviderError("initialize Claude ACP wrapper", err)
	}

	if err := c.createSession(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Connection) createSession(ctx context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get Claude ACP session working directory: %w", err)
	}

	session, err := c.conn.NewSession(ctx, acp.NewSessionRequest{
		Cwd:        cwd,
		McpServers: []acp.McpServer{},
	})
	if err != nil {
		return c.wrapProviderError("create Claude ACP session", err)
	}

	c.sessionID = session.SessionId
	return nil
}

func (c *Connection) Ask(_ string, _ chat.Conversation, text string) (string, error) {
	c.askMu.Lock()
	defer c.askMu.Unlock()

	if c.sessionID == "" {
		return "", errors.New("Claude ACP session is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), promptTimeout)
	defer cancel()

	c.response.Clear()
	_, err := c.conn.Prompt(ctx, acp.PromptRequest{
		SessionId: c.sessionID,
		Prompt:    []acp.ContentBlock{acp.TextBlock(text)},
	})
	if err != nil {
		return "", c.wrapProviderError("prompt Claude ACP wrapper", err)
	}

	var response strings.Builder
	for _, notification := range c.response.Snapshot() {
		c.appendAgentMessage(&response, c.sessionID, notification)
	}
	return response.String(), nil
}

func (c *Connection) Close() error {
	var closeErr error
	c.closeOnce.Do(func() {
		if c.stdin != nil {
			if err := c.stdin.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
				closeErr = fmt.Errorf("close Claude ACP stdin: %w", err)
			}
		}

		if c.cmd != nil && c.cmd.Process != nil && c.cmd.ProcessState == nil {
			select {
			case err := <-c.waitCh:
				if normalized := normalizeProcessWaitError(err); normalized != nil && closeErr == nil {
					closeErr = normalized
				}
				return
			case <-time.After(shutdownGracePeriod):
			}

			_ = c.cmd.Process.Kill()
		}

		if c.waitCh != nil {
			select {
			case <-c.waitCh:
			case <-time.After(shutdownKillWait):
				if closeErr == nil {
					closeErr = errors.New("timed out waiting for Claude ACP wrapper to exit")
				}
			}
		}
	})

	return closeErr
}

type streamingResponse struct {
	mu            sync.Mutex
	notifications []acp.SessionNotification
}

func (r *streamingResponse) Append(notification acp.SessionNotification) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notifications = append(r.notifications, notification)
}

func (r *streamingResponse) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notifications = nil
}

func (r *streamingResponse) Snapshot() []acp.SessionNotification {
	r.mu.Lock()
	defer r.mu.Unlock()

	return slices.Clone(r.notifications)
}

func (*Connection) appendAgentMessage(response *strings.Builder, sessionID acp.SessionId, notification acp.SessionNotification) {
	if notification.SessionId != sessionID {
		return
	}
	if notification.Update.AgentMessageChunk == nil {
		return
	}

	text := notification.Update.AgentMessageChunk.Content.Text
	if text == nil {
		return
	}

	response.WriteString(text.Text)
}

func (c *Connection) wrapProviderError(action string, err error) error {
	message := fmt.Sprintf("%s: %v", action, err)
	if stderr := strings.TrimSpace(c.stderr.String()); stderr != "" {
		message = fmt.Sprintf("%s (stderr: %s)", message, shorten(stderr, 400))
	}
	return errors.New(message)
}
