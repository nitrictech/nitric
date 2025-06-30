package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nitric",
	Long:  `Login to the Nitric CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.GetOrRefreshWorkosToken()
		if err == nil {
			user := fmt.Sprintf("%s %s <%s>", token.User.FirstName, token.User.LastName, token.User.Email)

			fmt.Printf("\n%s Already logged in as %s\n", style.Green(icons.Check), style.Teal(user))
			return
		}

		fmt.Printf("\n%s Logging in...\n", style.Purple(icons.Lightning+" Nitric"))

		err = auth.PerformPKCEFlow()
		if err != nil {
			fmt.Printf("\n%s Error logging in: %s\n", style.Red(icons.Cross), err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
