package main

import (
	"os"
	"strings"

	"github.com/nitrictech/nitric/cli/cmd"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/nitrictech/nitric/cli/pkg/cli"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func main() {

	injector := do.New()

	do.Provide(injector, func(inj *do.Injector) (workos.TokenStore, error) {
		config := do.MustInvoke[*config.Config](inj)
		apiUrl := config.GetNitricServerUrl()

		tokenStore, err := workos.NewKeyringTokenStore("nitric.v2.cli", apiUrl.String())
		if err != nil {
			return nil, err
		}
		return tokenStore, nil
	})

	do.Provide(injector, api.NewNitricApiClient)

	do.Provide(injector, func(inj *do.Injector) (*workos.WorkOSAuth, error) {
		tokenStore := do.MustInvoke[workos.TokenStore](inj)

		apiClient := do.MustInvoke[*api.NitricApiClient](inj)
		workosDetails, err := apiClient.GetWorkOSPublicDetails()
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset by peer") {
				cobra.CheckErr("failed to connect to the Nitric API. Please check your connection and try again. If the problem persists, please contact support.")
			}

			cobra.CheckErr(err)
		}
		return workos.NewWorkOSAuth(tokenStore, workosDetails.ClientID, workosDetails.ApiHostname), nil
	})

	do.Provide(injector, func(inj *do.Injector) (api.TokenProvider, error) {
		return do.Invoke[*workos.WorkOSAuth](inj)
	})

	do.Provide(injector, cli.NewCLI)
	do.Provide(injector, cli.NewAuthApp)

	rootCmd := cmd.NewRootCmd(injector)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
