package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd(injector do.Injector) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  `Manage the CLI configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			conf := do.MustInvoke[*config.Config](injector)

			if conf.FileUsed() != "" {
				fmt.Printf("file: %s\n\n", style.Teal(conf.FileUsed()))
				fmt.Println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, false, false, false).BorderForeground(colors.Teal).Render(conf.Dump()))
			} else {
				fmt.Println("file: none (using defaults)")
			}
		},
	}

	configSetCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set the CLI configuration",
		Long:  `Set the CLI configuration.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]

			conf := do.MustInvoke[*config.Config](injector)

			if err := conf.SetValue(key, value); err != nil {
				fmt.Printf("error setting config: %v", err)
			}

			if err := conf.Save(); err != nil {
				fmt.Printf("error saving config: %v", err)
			}

			fmt.Printf("Config set: %s: %s\n", key, value)
		},
	}

	configCmd.AddCommand(configSetCmd)
	return configCmd
}
