package cmd

import (
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Nitric",
	Long:  `Logout from the Nitric CLI.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := auth.DeleteWorkosToken()
		if err != nil && !errors.Is(err, auth.ErrNotFound) {
			fmt.Printf("\n%s Error logging out: %s\n", style.Red(icons.Cross), err)
			return
		}

		fmt.Printf("\n%s Logged out\n", style.Green(icons.Check))
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
