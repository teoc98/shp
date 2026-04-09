package queryinput

import (
	"container/list"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/teoc98/shp/tui/theme"
)

type Bubble struct {
	Styles    Styles
	textarea  textarea.Model

	history         *list.List
	historyMaxLen   int
	historySelected *list.Element
}

func New(jqtheme theme.Theme) Bubble {
	s := DefaultStyles()
	s.containerStyle.BorderForeground(jqtheme.Primary)
	ti := textarea.New()
	ti.Focus()
	// ti.FocusedStyle.Height(1)
	// ti.TextStyle.Height(1)
	ti.Placeholder = "rev | cowsay"
	ti.ShowLineNumbers = true

	promptLine := "$ "
	promptWidth := len(promptLine)
	promptFunc := func (lineIdx int) string {
		if lineIdx == 0 {
			return lipgloss.NewStyle().Bold(true).Foreground(jqtheme.Secondary).Render(promptLine)
		} else {
			return ""
		}
	}
	ti.SetPromptFunc(promptWidth, promptFunc)

	return Bubble{
		Styles:    s,
		textarea: ti,

		history:       list.New(),
		historyMaxLen: 512,
	}
}

func (b *Bubble) SetBorderColor(color lipgloss.TerminalColor) {
	b.Styles.containerStyle = b.Styles.containerStyle.BorderForeground(color)
}

func (b Bubble) GetInputValue() string {
	return b.textarea.Value()
}

func (b *Bubble) RotateHistory() {
	b.history.PushFront(b.textarea.Value())
	b.historySelected = b.history.Front()
	if b.history.Len() > b.historyMaxLen {
		b.history.Remove(b.history.Back())
	}
}

func (Bubble) Init() tea.Cmd {
	return textarea.Blink
}

func (b *Bubble) SetWidth(width int) {
	b.Styles.containerStyle = b.Styles.containerStyle.Width(width - b.Styles.containerStyle.GetHorizontalFrameSize())
	b.textarea.SetWidth(width - b.Styles.containerStyle.GetHorizontalFrameSize() - 1)
}

func (b Bubble) View() string {
	return b.Styles.containerStyle.Render(b.textarea.View())
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return b.updateKeyMsg(msg)
	default:
		var cmd tea.Cmd
		b.textarea, cmd = b.textarea.Update(msg)
		return b, cmd
	}
}

func (b *Bubble) SetQuery(query string) {
	b.textarea.SetValue(query)
}

func (b Bubble) updateKeyMsg(msg tea.KeyMsg) (Bubble, tea.Cmd) {
	switch msg.Type {
	// case tea.KeyUp:
	// 	return b.handleKeyUp()
	// case tea.KeyDown:
	// 	return b.handleKeyDown()
	// case tea.KeyEnter:
	// 	b.RotateHistory()
	// 	return b, nil
	default:
		var cmd tea.Cmd
		b.textarea, cmd = b.textarea.Update(msg)
		return b, cmd
	}
}

func (b Bubble) handleKeyUp() (Bubble, tea.Cmd) {
	if b.history.Len() == 0 {
		return b, nil
	}
	n := b.historySelected.Next()
	if n != nil {
		b.textarea.SetValue(n.Value.(string))
		b.textarea.CursorEnd()
		b.historySelected = n
	}
	return b, nil
}

func (b Bubble) handleKeyDown() (Bubble, tea.Cmd) {
	if b.history.Len() == 0 {
		return b, nil
	}
	p := b.historySelected.Prev()
	if p != nil {
		b.textarea.SetValue(p.Value.(string))
		b.textarea.CursorEnd()
		b.historySelected = p
	}
	return b, nil
}
