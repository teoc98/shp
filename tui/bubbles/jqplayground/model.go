package jqplayground

import (
	"os"
	"time"

	"github.com/teoc98/shp/tui/bubbles/fileselector"
	"github.com/teoc98/shp/tui/bubbles/help"
	"github.com/teoc98/shp/tui/bubbles/inputdata"
	"github.com/teoc98/shp/tui/bubbles/output"
	"github.com/teoc98/shp/tui/bubbles/queryinput"
	"github.com/teoc98/shp/tui/bubbles/state"
	"github.com/teoc98/shp/tui/bubbles/statusbar"
	"github.com/teoc98/shp/tui/theme"
)

func (b Bubble) GetState() state.State {
	return b.state
}

type Bubble struct {
	width            int
	height           int
	workingDirectory string
	state            state.State
	queryinput       queryinput.Bubble
	inputdata        inputdata.Bubble
	output           output.Bubble
	help             help.Bubble
	statusbar        statusbar.Bubble
	fileselector     fileselector.Bubble
	results          string
	cancel           func()
	theme            theme.Theme
	shell            string
	ExitMessage      string
	isJSONLines      bool
	showInputPanel   bool
}

func New(inputJSON []byte, filename string, query string, jqtheme theme.Theme, shell string) (Bubble, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return Bubble{}, err
	}

	sb := statusbar.New(jqtheme)
	sb.StatusMessageLifetime = time.Second * 10
	fs := fileselector.New(jqtheme)

	fs.SetInput(workingDirectory)

	inputData, err := inputdata.New(inputJSON, filename, jqtheme)
	if err != nil {
		return Bubble{}, err
	}
	queryInput := queryinput.New(jqtheme)
	if query != "" {
		queryInput.SetQuery(query)
	}

	b := Bubble{
		workingDirectory: workingDirectory,
		state:            state.Loading,
		queryinput:       queryInput,
		inputdata:        inputData,
		output:           output.New(jqtheme),
		help:             help.New(jqtheme),
		statusbar:        sb,
		fileselector:     fs,
		theme:            jqtheme,
		shell:            shell,
		showInputPanel:   true,
	}
	return b, nil
}
