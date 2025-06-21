package cmd

import (
	"github.com/nitrictech/nitric/cli/internal/config"
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

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nitric.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Load(cfgFile)
}
