package templates

import (
	"github.com/spf13/cobra"
)

var TemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage Nitric templates",
	Long:  `Manage Nitric templates.`,
}
