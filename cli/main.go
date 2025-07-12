package main

import (
	"os"

	"github.com/nitrictech/nitric/cli/cmd"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/config"
	details_service "github.com/nitrictech/nitric/cli/internal/details/service"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/nitrictech/nitric/cli/pkg/app"
	"github.com/samber/do/v2"
)

func createTokenStore(inj do.Injector) (*workos.KeyringTokenStore, error) {
	config := do.MustInvoke[*config.Config](inj)
	apiUrl := config.GetNitricServerUrl()

	tokenStore, err := workos.NewKeyringTokenStore("nitric.v2.cli", apiUrl.String())
	if err != nil {
		return nil, err
	}
	return tokenStore, nil
}

func main() {
	injector := do.New()

	do.Provide(injector, createTokenStore)
	do.Provide(injector, api.NewNitricApiClient)
	do.Provide(injector, details_service.NewService)
	do.Provide(injector, workos.NewWorkOSAuth)
	do.Provide(injector, app.NewNitricApp)
	do.Provide(injector, app.NewAuthApp)

	rootCmd := cmd.NewRootCmd(injector)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
