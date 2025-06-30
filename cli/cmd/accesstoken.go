package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/spf13/cobra"
)

var accessTokenCmd = &cobra.Command{
	Use:   "access-token",
	Short: "Print an access token for the Nitric Platform",
	Long:  `Print an access token for the Nitric Platform, using the current login session.`,
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
