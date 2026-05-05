package tui_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/bilus/experiments-chat/internal/buf"
	"github.com/bilus/experiments-chat/internal/chat"
	"github.com/bilus/experiments-chat/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
)

const (
	claudeReadyMessage   = "Claude is ready."
	programStartTimeout  = 50 * time.Millisecond
	programStopTimeout   = time.Second
	providerReadyTimeout = 30 * time.Second
)

type programResult struct {
	model tea.Model
	err   error
}

type testInput struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func newTestInput() *testInput {
	reader, writer := io.Pipe()
	return &testInput{
		reader: reader,
		writer: writer,
	}
}

func (i *testInput) close() {
	_ = i.writer.Close()
	_ = i.reader.Close()
}

type terminalChatFeature struct {
	program *tea.Program
	done    chan programResult
	input   *testInput
	model   tui.Model
	output  buf.Buffer
}

func (f *terminalChatFeature) reset() {
	f.cleanup()
	*f = terminalChatFeature{}
}

func (f *terminalChatFeature) cleanup() {
	if f.program == nil || f.done == nil {
		return
	}

	program := f.program
	done := f.done
	input := f.input

	f.quitProgram()
	select {
	case result := <-done:
		_ = f.storeProgramResult(result)
		return
	case <-time.After(programStopTimeout):
	}

	program.Kill()
	if input != nil {
		input.close()
	}

	select {
	case <-done:
	case <-time.After(programStopTimeout):
	}

	f.clearProgram()
}

func (f *terminalChatFeature) clearProgram() {
	f.program = nil
	f.done = nil
	f.input = nil
}

func (f *terminalChatFeature) storeProgramResult(result programResult) error {
	if f.input != nil {
		f.input.close()
	}
	defer f.clearProgram()

	if result.err != nil {
		return result.err
	}

	model, ok := result.model.(tui.Model)
	if !ok {
		return fmt.Errorf("expected terminal chat model, got %T", result.model)
	}

	f.model = model
	if err := f.model.Close(); err != nil {
		return fmt.Errorf("close terminal chat model: %w", err)
	}
	return nil
}

func (f *terminalChatFeature) captureFinalModel(timeoutMessage string) error {
	if f.done == nil {
		return fmt.Errorf("terminal chat application is not running")
	}

	select {
	case result := <-f.done:
		if err := f.storeProgramResult(result); err != nil {
			return fmt.Errorf("run terminal chat: %w", err)
		}

		return nil
	case <-time.After(programStopTimeout):
		return errors.New(timeoutMessage)
	}
}

func (f *terminalChatFeature) quitProgram() {
	if f.program == nil {
		return
	}

	f.program.Send(tea.KeyMsg{Type: tea.KeyCtrlD})
}

func (f *terminalChatFeature) renderedView() (string, error) {
	if f.program != nil {
		f.quitProgram()
		if err := f.captureFinalModel("terminal chat did not exit while capturing rendered output"); err != nil {
			return "", err
		}
	}

	return f.model.View(), nil
}

func (f *terminalChatFeature) renderedOutput() (string, error) {
	if f.program != nil {
		f.quitProgram()
		if err := f.captureFinalModel("terminal chat did not exit while capturing rendered output"); err != nil {
			return "", err
		}
	}

	return f.output.String(), nil
}

func (f *terminalChatFeature) waitForOutput(expected string) error {
	deadline := time.After(providerReadyTimeout)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	for {
		if strings.Contains(f.output.String(), expected) {
			return nil
		}

		select {
		case result := <-f.done:
			if err := f.storeProgramResult(result); err != nil {
				return fmt.Errorf("terminal chat exited while waiting for %q: %w", expected, err)
			}

			return fmt.Errorf("terminal chat exited before rendering %q; output: %q", expected, f.output.String())
		case <-deadline:
			return fmt.Errorf("timed out waiting for %q; output: %q", expected, f.output.String())
		case <-tick.C:
		}
	}
}

func (f *terminalChatFeature) theTerminalChatApplicationIsRunning() error {
	provider := chat.StaticProvider{Response: "I don't understand."}
	return f.startTerminalChat(tui.NewModel(chat.NewAgent("", provider)))
}

func (f *terminalChatFeature) claudeIsTheSelectedProvider() error {
	f.model = tui.NewWithClaude()
	return nil
}

func (f *terminalChatFeature) theTerminalChatApplicationStarts() error {
	return f.startTerminalChat(f.model)
}

func (f *terminalChatFeature) startTerminalChat(model tui.Model) error {
	f.model = model
	input := newTestInput()
	program := tea.NewProgram(f.model, tea.WithInput(input.reader), tea.WithOutput(&f.output))
	done := make(chan programResult, 1)

	f.program = program
	f.done = done
	f.input = input

	go func() {
		model, err := program.Run()
		done <- programResult{model: model, err: err}
	}()

	select {
	case result := <-done:
		if err := f.storeProgramResult(result); err != nil {
			return fmt.Errorf("terminal chat exited while starting: %w", err)
		}

		return fmt.Errorf("terminal chat exited before accepting input")
	case <-time.After(programStartTimeout):
		return nil
	}
}

func (f *terminalChatFeature) connectionToClaudeIsSuccessful() error {
	return f.waitForOutput(claudeReadyMessage)
}

func (f *terminalChatFeature) theSystemShouldIndicateReadiness() error {
	output, err := f.renderedOutput()
	if err != nil {
		return err
	}

	if !strings.Contains(output, claudeReadyMessage) {
		return fmt.Errorf("expected terminal chat output to contain %q, got %q", claudeReadyMessage, output)
	}

	return nil
}

func (f *terminalChatFeature) theTerminalChatShouldShow(expected string) error {
	view, err := f.renderedView()
	if err != nil {
		return err
	}

	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected terminal chat to show %q, got %q", expected, view)
	}

	return nil
}

func (f *terminalChatFeature) sendKey(key tea.KeyType) error {
	if f.program == nil {
		return fmt.Errorf("terminal chat application is not running")
	}

	f.program.Send(tea.KeyMsg{Type: key})
	return nil
}

func (f *terminalChatFeature) sendText(text string) error {
	if f.program == nil {
		return fmt.Errorf("terminal chat application is not running")
	}

	f.program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)})
	return nil
}

func (f *terminalChatFeature) theUserTypes(text string) error {
	return f.sendText(text)
}

func (f *terminalChatFeature) theUserPressesLeftArrow() error {
	return f.sendKey(tea.KeyLeft)
}

func (f *terminalChatFeature) theUserPressesBackspace() error {
	return f.sendKey(tea.KeyBackspace)
}

func (f *terminalChatFeature) theUserSubmitsThePrompt() error {
	return f.sendKey(tea.KeyEnter)
}

func (f *terminalChatFeature) theUserPressesCtrlD() error {
	return f.sendKey(tea.KeyCtrlD)
}

func (f *terminalChatFeature) theTerminalChatShouldExit() error {
	return f.captureFinalModel("terminal chat did not exit after Ctrl+D")
}

func (f *terminalChatFeature) theTerminalChatShouldRender(expected string) error {
	view, err := f.renderedView()
	if err != nil {
		return err
	}

	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected terminal chat to render %q, got %q", expected, view)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			feature := &terminalChatFeature{}

			sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				feature.reset()
				return ctx, nil
			})
			sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				feature.cleanup()
				return ctx, nil
			})

			sc.Step(`^the terminal chat application is running$`, feature.theTerminalChatApplicationIsRunning)
			sc.Step(`^Claude is the selected provider$`, feature.claudeIsTheSelectedProvider)
			sc.Step(`^the terminal chat application starts$`, feature.theTerminalChatApplicationStarts)
			sc.Step(`^connection to Claude is successful$`, feature.connectionToClaudeIsSuccessful)
			sc.Step(`^the system should indicate readiness$`, feature.theSystemShouldIndicateReadiness)
			sc.Step(`^the user types "([^"]*)"$`, feature.theUserTypes)
			sc.Step(`^the user presses Left Arrow$`, feature.theUserPressesLeftArrow)
			sc.Step(`^the user presses Backspace$`, feature.theUserPressesBackspace)
			sc.Step(`^the user presses Ctrl\+D$`, feature.theUserPressesCtrlD)
			sc.Step(`^the user submits the prompt$`, feature.theUserSubmitsThePrompt)
			sc.Step(`^the terminal chat should show "([^"]*)"$`, feature.theTerminalChatShouldShow)
			sc.Step(`^the terminal chat should render "([^"]*)"$`, feature.theTerminalChatShouldRender)
			sc.Step(`^the terminal chat should exit$`, feature.theTerminalChatShouldExit)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("godog feature suite failed")
	}
}
