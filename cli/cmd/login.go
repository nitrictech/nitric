package cmd

import (
	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nitric",
	Long:  `Login to Nitric.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkce := auth.WorkOsPKCE{}
		pkce.PerformPKCEFlow()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
