package jqplayground

import (
	"context"
	"io"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/noahgorstein/jqp/tui/utils"
)

type errorMsg struct {
	error error
}

type queryResultMsg struct {
	rawResults         string
	highlightedResults string
}

type writeToFileMsg struct{}

type copyQueryToClipboardMsg struct{}

type copyResultsToClipboardMsg struct{}

// processQueryResults iterates through the results of a gojq query on the provided JSON object
// and appends the formatted results to the provided string builder.
func processQueryResults(ctx context.Context, results *strings.Builder, shell string, script string, data []byte) error {
	cmd := exec.Command(shell, "-c", script)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(data)
	}()

	cmdout, err := io.ReadAll(stdout)
	if err != nil {
		return err
	}
	cmderr, err := io.ReadAll(stderr)
	if err != nil {
		return err
	}
	fmt.Fprintf(results, "%s\n", cmdout)
	fmt.Fprintf(results, "%s\n", cmderr)

	err = cmd.Wait()
	return err
}

func processJSONWithQuery(ctx context.Context, results *strings.Builder, shell string, script string, data []byte) error {
	err := processQueryResults(ctx, results, shell, script, data)
	if err != nil {
		return err
	}

	return nil
}

func processJSONLinesWithQuery(ctx context.Context, results *strings.Builder, shell string, script string, data []byte) error {
	const maxBufferSize = 100 * 1024 * 1024 // 100MB max buffer size

	processLine := func(line []byte) error {
		return processJSONWithQuery(ctx, results, shell, script, line)
	}

	return utils.ScanLinesWithDynamicBufferSize(data, maxBufferSize, processLine)
}

func (b *Bubble) executeQueryOnInput(ctx context.Context) (string, error) {
	var results strings.Builder
	shell := b.shell
	script := b.queryinput.GetInputValue()

	processor := processJSONWithQuery

	if b.isJSONLines {
		processor = processJSONLinesWithQuery
	}
	err := processor(ctx, &results, shell, script, b.inputdata.GetInputJSON())
	return results.String(), err
}

func (b *Bubble) executeQueryCommand(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		results, err := b.executeQueryOnInput(ctx)
		if err != nil {
			return errorMsg{error: err}
		}
		highlightedOutput, err := utils.Prettify([]byte(results), b.theme.ChromaStyle, true)
		if err != nil {
			return errorMsg{error: err}
		}
		return queryResultMsg{
			rawResults:         results,
			highlightedResults: highlightedOutput.String(),
		}
	}
}

func (b Bubble) saveOutput() tea.Cmd {
	if b.fileselector.GetInput() == "" {
		return b.copyOutputToClipboard()
	}
	return b.writeOutputToFile()
}

func (b Bubble) copyOutputToClipboard() tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(b.results)
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return copyResultsToClipboardMsg{}
	}
}

func (b Bubble) writeOutputToFile() tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(b.fileselector.GetInput(), []byte(b.results), 0o600)
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return writeToFileMsg{}
	}
}

func (b Bubble) copyQueryToClipboard() tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(b.queryinput.GetInputValue())
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return copyQueryToClipboardMsg{}
	}
}
