package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

var (
	debugFlag bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nitric",
	Long:  `Login to the Nitric CLI.`,
	Run: func(cmd *cobra.Command, args []string) {

		token, err := auth.GetOrRefreshWorkosToken()
		if err == nil {

			if debugFlag {
				tokenJson, err := json.MarshalIndent(token, "", "  ")
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(tokenJson))
			}

			user := fmt.Sprintf("%s %s <%s>", token.User.FirstName, token.User.LastName, token.User.Email)

			fmt.Printf("\n%s Already logged in as %s\n", style.Green(icons.Check), style.Teal(user))
			return
		}

		if debugFlag {
			fmt.Printf("Error getting workos token: %v\n", err)
		}

		pkce, err := auth.NewWorkOsPKCE()
		cobra.CheckErr(err)

		err = pkce.PerformPKCEFlow()
		if err != nil {
			fmt.Printf("\n%s Error logging in: %s\n", style.Red(icons.Cross), err)
			return
		}
	},
}

func init() {
	loginCmd.Flags().BoolVar(&debugFlag, "debug", false, "Debug mode")
	rootCmd.AddCommand(loginCmd)
}
