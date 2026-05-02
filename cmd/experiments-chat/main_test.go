package main

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/bilus/experiments-chat/internal/buf"
	"github.com/bilus/experiments-chat/internal/test"
)

const claudeCommandOverrideEnvName = "EXPERIMENTS_CHAT_CLAUDE_ACP_COMMAND"

func TestFakeAgent(t *testing.T) {
	test.RunFakeAgent(t)
}

func TestRunStartsClaudeProviderFromFlag(t *testing.T) {
	test.UseFakeAgent(t, claudeCommandOverrideEnvName)

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	output := &buf.Buffer{}
	done := make(chan error, 1)
	go func() {
		done <- run([]string{"--provider", "claude"}, reader, output, io.Discard)
	}()

	waitForOutput(t, output, "Claude is ready.")
	if _, err := writer.Write([]byte("hello\r")); err != nil {
		t.Fatalf("submit prompt: %v", err)
	}
	waitForOutput(t, output, "helper heard: hello")

	if _, err := writer.Write([]byte{4}); err != nil {
		t.Fatalf("send Ctrl+D: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run terminal chat: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("terminal chat did not exit after Ctrl+D")
	}
}

func waitForOutput(t *testing.T, output *buf.Buffer, expected string) {
	t.Helper()

	deadline := time.After(2 * time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	for {
		if strings.Contains(output.String(), expected) {
			return
		}

		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %q; output: %q", expected, output.String())
		case <-tick.C:
		}
	}
}
