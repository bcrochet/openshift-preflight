// Package registration contains all plugins to register.
// Plugin developers should blank-initialize their plugins here
package registration

// Plugin initialization
import (
	_ "github.com/opdev/container-certification"
	_ "github.com/opdev/container-certification/rootexception"
	_ "github.com/opdev/container-certification/scratchexception"
	_ "github.com/opdev/operator-certification"
	// _ "github.com/opdev/plugin-template"
)
