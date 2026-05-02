package claude

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandPrefersClaudeAgentACP(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, "claude"))
	wrapperPath := filepath.Join(binDir, "claude-agent-acp")
	writeExecutable(t, wrapperPath)

	t.Setenv(commandOverrideEnvName, "")
	t.Setenv("PATH", binDir)

	command, err := command()
	if err != nil {
		t.Fatalf("resolve Claude ACP command: %v", err)
	}

	if command[0] != wrapperPath {
		t.Fatalf("expected claude-agent-acp command %q, got %q", wrapperPath, command[0])
	}
}

func TestCommandDoesNotRequireClaudeCLI(t *testing.T) {
	binDir := t.TempDir()
	wrapperPath := filepath.Join(binDir, "claude-agent-acp")
	writeExecutable(t, wrapperPath)

	t.Setenv(commandOverrideEnvName, "")
	t.Setenv("PATH", binDir)

	command, err := command()
	if err != nil {
		t.Fatalf("resolve Claude ACP command: %v", err)
	}

	if command[0] != wrapperPath {
		t.Fatalf("expected claude-agent-acp command %q, got %q", wrapperPath, command[0])
	}
}

func writeExecutable(t *testing.T, path string) {
	t.Helper()

	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", path, err)
	}
}
