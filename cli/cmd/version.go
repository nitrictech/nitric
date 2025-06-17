package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/spf13/cobra"
)

var highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Long:  `Display the version number of the Nitric CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nitric cli version %s\n", highlight.Render(version.Version))

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
