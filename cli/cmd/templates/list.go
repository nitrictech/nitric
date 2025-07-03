package templates

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Nitric templates",
	Long:  `List available Nitric templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewNitricApiClient(config.GetNitricServerUrl(), auth.WithAuthHeader)

		resp, err := client.GetTemplates()
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

func init() {
	TemplatesCmd.AddCommand(listCmd)
}
