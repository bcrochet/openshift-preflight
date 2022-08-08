package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/prompt"
)

type OperatorPromptModel struct {
	inputs     []*prompt.Model
	focusIndex int
	done       bool
	err        error
}

func NewOperatorPrompt() *OperatorPromptModel {
	m := &OperatorPromptModel{
		inputs: make([]*prompt.Model, 2),
		err:    nil,
	}

	var t prompt.Model
	for i := range m.inputs {

		switch i {
		case 0:
			t = prompt.New("Project ID?")
			t.Template = customTemplate
			t.ResultTemplate = customResultTemplate
			t.CharLimit = 16
			t.Placeholder = "Project ID required"
		case 1:
			t = prompt.New("Pull Request URL?")
			t.Template = customTemplate
			t.ResultTemplate = customResultTemplate
			t.CharLimit = 16
			t.Placeholder = "Pull Request URL, https://..."
			t.CharLimit = 255
		}

		m.inputs[i] = t
	}

	return m
}

func (c OperatorPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (c *OperatorPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return c, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && c.focusIndex == len(c.inputs) {
				return c, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			if c.focusIndex > len(c.inputs) {
				c.focusIndex = 0
			} else if c.focusIndex < 0 {
				c.focusIndex = len(c.inputs)
			}

			cmds := make([]tea.Cmd, len(c.inputs))
			for i := 0; i <= len(c.inputs)-1; i++ {
				if i == c.focusIndex {
					// Set focused state
					continue
				}
				// Remove focused state
			}

			return c, tea.Batch(cmds...)
		}
	}

	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *OperatorPromptModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c OperatorPromptModel) View() string {
	var b strings.Builder

	for i := range c.inputs {
		b.WriteString(c.inputs[i].View())
		if i < len(c.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if c.focusIndex == len(c.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
