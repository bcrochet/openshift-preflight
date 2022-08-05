package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/tui"
	"github.com/spf13/cobra"
)

func supportOperatorCmd() *cobra.Command {
	supportOperator := &cobra.Command{
		Use:   "operator <your project ID> <pullRequestURL>",
		Short: "Creates a support request for an operator",
		Args: func(cmd *cobra.Command, args []string) error {
			if !interactive && len(args) < 2 {
				return fmt.Errorf("a project ID and pull request URL are required")
			}
			return nil
		},
		Long: `Generate a URL that can be used to open a ticket with Red Hat Support if you're having an issue passing certification checks.`,
		RunE: supportOperatorRunE,
	}

	return supportOperator
}

func supportOperatorRunE(cmd *cobra.Command, args []string) error {
	if interactive {
		if err := tea.NewProgram(tui.NewOperatorPrompt()).Start(); err != nil {
			return err
		}
		return nil
	}

	ptype := "operator"
	pid := args[0]
	prurl := args[1]

	support, err := newSupportTextGenerator(ptype, pid, prurl)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), support.Generate())
	return nil
}
