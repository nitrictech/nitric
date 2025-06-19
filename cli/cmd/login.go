package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/auth/token"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nitric",
	Long:  `Login to the Nitric CLI.`,
	Run: func(cmd *cobra.Command, args []string) {

		token, err := token.GetWorkosToken()
		if err == nil {
			user := fmt.Sprintf("%s %s <%s>", token.User.FirstName, token.User.LastName, token.User.Email)

			fmt.Printf("\n%s Already logged in as %s\n", style.Green(icons.Check), style.Teal(user))
			return
		}

		pkce := auth.WorkOsPKCE{}
		pkce.PerformPKCEFlow()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
