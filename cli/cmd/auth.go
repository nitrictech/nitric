package cmd

import (
	"github.com/nitrictech/nitric/cli/pkg/cli"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

// NewLoginCmd creates the login command
func NewLoginCmd(injector *do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to Nitric",
		Long:  `Login to the Nitric CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			app := do.MustInvoke[*cli.AuthApp](injector)
			app.Login()
		},
	}
}

// NewLogoutCmd creates the logout command
func NewLogoutCmd(injector *do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from Nitric",
		Long:  `Logout from the Nitric CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			app := do.MustInvoke[*cli.AuthApp](injector)
			app.Logout()
		},
	}
}

// NewAccessTokenCmd creates the access token command
func NewAccessTokenCmd(injector *do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "access-token",
		Short: "Get access token",
		Long:  `Get the current access token.`,
		Run: func(cmd *cobra.Command, args []string) {
			app := do.MustInvoke[*cli.AuthApp](injector)
			app.AccessToken()
		},
	}
}
