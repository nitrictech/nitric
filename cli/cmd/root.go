package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func NewRootCmd(injector do.Injector) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nitric",
		Short: fmt.Sprintf("%s CLI - The command line interface for %s", version.ProductName, version.ProductName),
		Long:  fmt.Sprintf("%s CLI is a command line interface for managing and deploying %s applications. It provides a set of commands to help you develop, test, and deploy your %s applications.", version.ProductName, version.ProductName, version.ProductName),
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
	rootCmd.AddCommand(NewTeamCmd(injector))

	return rootCmd
}

// NewVersionCmd creates the version command
func NewVersionCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the %s CLI version", version.ProductName),
		Long:  fmt.Sprintf("Display the version number of the %s CLI.", version.ProductName),
		Run: func(cmd *cobra.Command, args []string) {
			highlight := lipgloss.NewStyle().Foreground(colors.Teal)
			fmt.Printf("%s cli version %s\n", version.ProductName, highlight.Render(version.Version))
		},
	}
}
