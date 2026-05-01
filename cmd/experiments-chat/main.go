package main

import (
	"fmt"
	"os"

	"github.com/bilus/experiments-chat/internal/chat"
	"github.com/bilus/experiments-chat/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	agent := chat.NewAgent("", chat.StaticProvider{Response: "I don't understand."})

	if _, err := tea.NewProgram(tui.NewModel(agent)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
