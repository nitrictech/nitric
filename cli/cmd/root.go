package cmd

import (
	"log"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "nitric",
		Short: "Nitric CLI - The command line interface for Nitric",
		Long: `Nitric CLI is a command line interface for managing and deploying
Nitric applications. It provides a set of commands to help you develop,
test, and deploy your Nitric applications.`,
	}
)

type Dependencies struct {
	WorkOSAuth      *workos.WorkOSAuth
	NitricApiClient *api.NitricApiClient
}

var deps *Dependencies = &Dependencies{}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initDependencies)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nitric.yaml)")

	rootCmd.AddCommand(NewLoginCmd(deps))
	rootCmd.AddCommand(NewLogoutCmd(deps))
	rootCmd.AddCommand(NewAccessTokenCmd(deps))

	rootCmd.AddCommand(NewTemplatesCmd(deps))
	rootCmd.AddCommand(NewBuildCmd(deps))
}

func initDependencies() {
	deps.NitricApiClient = api.NewNitricApiClient(config.GetNitricServerUrl())

	workosDetails, err := deps.NitricApiClient.GetWorkOSPublicDetails()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset by peer") {
			log.Fatal("unable to connect to the Nitric API. Please check your connection and try again. If the problem persists, please contact support.")
		}

		log.Fatal(err)
	}

	deps.WorkOSAuth = workos.NewWorkOSAuth(workos.NewKeyringTokenStore("nitric.v2.cli"), workosDetails.ClientID, workosDetails.ApiHostname)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Load(cfgFile)
}
