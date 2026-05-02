package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bilus/experiments-chat/internal/chat"
	"github.com/bilus/experiments-chat/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, input io.Reader, output io.Writer, stderr io.Writer) error {
	flags := flag.NewFlagSet("experiments-chat", flag.ContinueOnError)
	flags.SetOutput(stderr)
	providerName := flags.String("provider", "static", "chat provider: static or claude")
	if err := flags.Parse(args); err != nil {
		return err
	}

	model, err := modelForProvider(*providerName)
	if err != nil {
		return err
	}

	finalModel, err := tea.NewProgram(model, tea.WithInput(input), tea.WithOutput(output)).Run()
	return errors.Join(err, closeModel(finalModel))
}

func modelForProvider(providerName string) (tui.Model, error) {
	switch providerName {
	case "static":
		provider := chat.StaticProvider{Response: "I don't understand."}
		return tui.NewModel(chat.NewAgent("", provider)), nil
	case "claude":
		return tui.NewWithClaude(), nil
	default:
		return tui.Model{}, fmt.Errorf("unknown provider %q", providerName)
	}
}

func closeModel(model tea.Model) error {
	tuiModel, ok := model.(tui.Model)
	if !ok {
		return nil
	}

	return tuiModel.Close()
}
