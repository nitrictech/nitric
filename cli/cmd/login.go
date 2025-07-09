package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

func NewLoginCmd(deps *Dependencies) *cobra.Command {
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to Nitric",
		Long:  `Login to the Nitric CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			user, err := deps.WorkOSAuth.Login()
			if err != nil {
				fmt.Printf("\n%s Error logging in: %s\n", style.Red(icons.Cross), err)
				return
			}

			fmt.Printf("\n%s Login successful, welcome %s\n", style.Green(icons.Check), style.Teal(user.FirstName))
		},
	}

	return loginCmd
}
