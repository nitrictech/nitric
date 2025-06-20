package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set the CLI configuration",
	Long:  `Set the CLI configuration.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := config.SetValue(key, value); err != nil {
			fmt.Printf("Error setting config: %v\n", err)
		}

		if err := config.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
		}

		fmt.Printf("Config set: %s: %s\n", key, value)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage the CLI configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() != "" {
			fmt.Printf("using config file: %s\n", style.Teal(viper.ConfigFileUsed()))
			items := config.GetAllConfigItems()

			if len(items) > 0 {
				fmt.Println("")
			}

			for key, value := range items {
				fmt.Printf("%s: %s\n", style.Teal(key), value)
			}
		} else {
			fmt.Println("Config file: none (using defaults)")
		}
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
