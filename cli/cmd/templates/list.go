package templates

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Nitric templates",
	Long:  `List available Nitric templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.GetOrRefreshWorkosToken()
		if err != nil {
			fmt.Printf("\n Not currently logged in, run `nitric login` to login")
			return
		}

		client := api.NewNitricApiClient(config.GetNitricServerUrl(), &token.AccessToken)

		resp, err := client.GetTemplates()
		if err != nil {
			fmt.Printf("Failed to get templates: %v", err)
			return
		}

		for _, template := range resp.Templates {
			fmt.Printf(" - %s\n", template.Name)
		}
	},
}

func init() {
	TemplatesCmd.AddCommand(listCmd)
}
