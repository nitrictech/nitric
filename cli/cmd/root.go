package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func NewRootCmd(injector do.Injector) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nitric",
		Short: "Nitric CLI - The command line interface for Nitric",
		Long: `Nitric CLI is a command line interface for managing and deploying
Nitric applications. It provides a set of commands to help you develop,
test, and deploy your Nitric applications.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.Load(cmd)
			if err != nil {
				return err
			}

			do.ProvideValue(injector, conf)
			return nil
		},
	}

	// Add commands that use the CLI struct methods
	rootCmd.AddCommand(NewLoginCmd(injector))
	rootCmd.AddCommand(NewLogoutCmd(injector))
	rootCmd.AddCommand(NewAccessTokenCmd(injector))
	rootCmd.AddCommand(NewVersionCmd(injector))
	rootCmd.AddCommand(NewTemplatesCmd(injector))
	rootCmd.AddCommand(NewInitCmd(injector))
	rootCmd.AddCommand(NewNewCmd(injector))
	rootCmd.AddCommand(NewBuildCmd(injector))
	rootCmd.AddCommand(NewGenerateCmd(injector))
	rootCmd.AddCommand(NewEditCmd(injector))
	rootCmd.AddCommand(NewDevCmd(injector))
	rootCmd.AddCommand(NewConfigCmd(injector))

	return rootCmd
}

// NewVersionCmd creates the version command
func NewVersionCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Long:  `Display the version number of the Nitric CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
			fmt.Printf("nitric cli version %s\n", highlight.Render(version.Version))
		},
	}
}
