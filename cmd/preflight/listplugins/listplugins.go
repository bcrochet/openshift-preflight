package listplugins

import (
	"fmt"

	"github.com/redhat-openshift-ecosystem/preflight"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-plugins",
		Short: "list the plugins that have been registered",
		RunE:  func(cmd *cobra.Command, args []string) error { return listPlugins() },
	}

	return cmd
}

func listPlugins() error {
	fmt.Println("listing plugins")
	registered := preflight.RegisteredPlugins()
	for k, v := range registered {
		fmt.Printf("Plugin '%s' at version %s is registered at entrypoint '%s'\n", v.Name(), v.Version(), k)
	}

	return nil
}
