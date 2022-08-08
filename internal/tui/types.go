package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type (
	tickMsg struct{}
	errMsg  error
)

const (
	projectIdKey   = "Project ID"
	pullRequestUrl = "Pull Request URL"

	customTemplate = `
	{{- "┏" }}━{{ Repeat "━" (Len .Prompt) }}━┯━{{ Repeat "━" 13 }}{{ "━━━━┓\n" }}
	{{- "┃" }} {{ Bold .Prompt }} │ {{ .Input -}}
	{{- Repeat " " (Max 0 (Sub 16 (Len .Input))) }}
	{{- if .ValidationError -}}
		{{- Foreground "1" (Bold "✘") -}}
	{{- else -}}
		{{- Foreground "2" (Bold "✔") -}}
	{{- end -}}┃
	{{- if .ValidationError -}}
		{{- (print " Error: " (Foreground "1" .ValidationError.Error)) -}}
	{{- end -}}
	{{- "\n┗" }}━{{ Repeat "━" (Len .Prompt) }}━┷━{{ Repeat "━" 13 }}{{ "━━━━┛\n" -}}
	{{- if .AutoCompleteIndecisive -}}
		{{ print "  Suggestions: " }}
		{{- range $suggestion := AutoCompleteSuggestions -}}
			{{- print $suggestion " " -}}
		{{- end -}}
	{{- end -}}
	`

	customResultTemplate = `
	Create a support ticket by:
		1. Copying URL: {{ .SupportUrl }}
		2. Paste above URL in a browser
		3. Login with Red Hat SSO
		4. Enter an issue summary and description
		5. Preview and Submit your ticket
	`
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
