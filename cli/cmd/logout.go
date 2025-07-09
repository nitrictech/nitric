package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

func NewLogoutCmd(deps *Dependencies) *cobra.Command {
	var logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Logout from Nitric",
		Long:  `Logout from the Nitric CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := deps.WorkOSAuth.Logout()
			if err != nil {
				fmt.Printf("\n%s Error logging out: %s\n", style.Red(icons.Cross), err)
				return
			}

			fmt.Printf("\n%s Logged out\n", style.Green(icons.Check))
		},
	}

	return logoutCmd
}
