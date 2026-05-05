package claude

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func command() ([]string, error) {
	if override := strings.TrimSpace(os.Getenv(commandOverrideEnvName)); override != "" {
		command, err := commandFromOverride(override)
		if err != nil {
			return nil, err
		}
		return resolveCommand(command)
	}

	if path, err := exec.LookPath("claude-agent-acp"); err == nil {
		return []string{path}, nil
	}

	return nil, errors.New("Claude ACP wrapper not found: install claude-agent-acp")
}

func commandFromOverride(override string) ([]string, error) {
	if strings.HasPrefix(override, "[") {
		var command []string
		if err := json.Unmarshal([]byte(override), &command); err != nil {
			return nil, fmt.Errorf("parse %s as JSON command array: %w", commandOverrideEnvName, err)
		}
		return command, nil
	}

	return []string{override}, nil
}

func resolveCommand(command []string) ([]string, error) {
	if len(command) == 0 {
		return nil, errors.New("command cannot be empty")
	}

	path, err := exec.LookPath(command[0])
	if err != nil {
		return nil, err
	}

	return append([]string{path}, command[1:]...), nil
}
