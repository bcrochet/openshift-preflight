package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ContainerPromptModel struct {
	inputs     []textinput.Model
	focusIndex int
	done       bool
	err        error
}

func NewContainerPrompt() *ContainerPromptModel {
	pi := textinput.New()
	pi.Placeholder = "Project ID required"
	pi.Focus()
	pi.PromptStyle = focusedStyle
	pi.CharLimit = 16

	return &ContainerPromptModel{
		inputs: []textinput.Model{
			pi,
		},
		err: nil,
	}
}

func (c ContainerPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (c ContainerPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			s := msg.String()

			if s == "enter" && c.focusIndex == len(c.inputs) {
				c.done = true
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
					cmds[i] = c.inputs[i].Focus()
					c.inputs[i].PromptStyle = focusedStyle
					c.inputs[i].TextStyle = focusedStyle
					continue
				}

				c.inputs[i].Blur()
				c.inputs[i].PromptStyle = noStyle
				c.inputs[i].TextStyle = noStyle
			}

			return c, tea.Batch(cmds...)
		}

	// We handle errors just like any other message
	case errMsg:
		c.err = msg
		return c, nil
	}

	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *ContainerPromptModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c ContainerPromptModel) View() string {
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
