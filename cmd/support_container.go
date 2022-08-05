package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/tui"
	"github.com/spf13/cobra"
)

func supportContainerCmd() *cobra.Command {
	supportContainer := &cobra.Command{
		Use:   "container <your project ID>",
		Short: "Creates a support request",
		Args: func(cmd *cobra.Command, args []string) error {
			if !interactive && len(args) < 1 {
				return fmt.Errorf("project ID not provided")
			}
			return nil
		},
		Long: `Generate a URL that can be used to open a ticket with Red Hat Support if you're having an issue passing certification checks.`,
		RunE: supportContainerRunE,
	}

	return supportContainer
}

func supportContainerRunE(cmd *cobra.Command, args []string) error {
	if interactive {
		if err := tea.NewProgram(tui.NewContainerPrompt()).Start(); err != nil {
			return err
		}
		return nil
	}

	// prurl is not needed for container support
	prurl := ""
	ptype := "container"
	pid := args[0]

	support, err := newSupportTextGenerator(ptype, pid, prurl)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), support.Generate())
	return nil
}
