package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAccessTokenCmd(deps *Dependencies) *cobra.Command {
	var accessTokenCmd = &cobra.Command{
		Use:   "access-token",
		Short: "Print an access token for the Nitric Platform",
		Long:  `Print an access token for the Nitric Platform, using the current login session.`,
		Run: func(cmd *cobra.Command, args []string) {
			token, err := deps.WorkOSAuth.GetAccessToken()
			if err != nil {
				fmt.Printf("\n Not currently logged in, run `nitric login` to login")
				return
			}

			fmt.Println(token)
		},
	}

	return accessTokenCmd
}
