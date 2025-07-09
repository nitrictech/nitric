package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/spf13/cobra"
)

func NewTemplatesCmd(deps *Dependencies) *cobra.Command {
	var templatesCmd = &cobra.Command{
		Use:   "templates",
		Short: "Manage Nitric templates",
		Long:  `Manage Nitric templates.`,
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available Nitric templates",
		Long:  `List available Nitric templates.`,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := deps.NitricApiClient.GetTemplates()
			if err != nil {
				fmt.Printf("Failed to get templates: %v", err)
				return
			}

			if len(resp) == 0 {
				fmt.Println("No templates found")
				return
			}

			templateStyle := lipgloss.NewStyle().Foreground(colors.Purple)

			fmt.Println(templateStyle.Render("\nAvailable templates:"))

			for _, template := range resp {
				fmt.Printf(" %s\n", template)
			}
		},
	}

	templatesCmd.AddCommand(listCmd)

	return templatesCmd
}
