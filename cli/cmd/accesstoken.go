package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/spf13/cobra"
)

var accessTokenCmd = &cobra.Command{
	Use:   "accesstoken",
	Short: "Get the access token for the Nitric Platform",
	Long:  `Get the access token for the Nitric Platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.GetOrRefreshWorkosToken()
		if err != nil {
			fmt.Printf("\n Not currently logged in, run `nitric login` to login")
			return
		}

		fmt.Println(token.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(accessTokenCmd)
}
