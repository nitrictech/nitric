package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd(injector do.Injector) *cobra.Command {
	var listFlag bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long: `Manage the CLI configuration.

By default, this command displays the current configuration values.
Use the --list flag to see all available configuration keys and their descriptions.`,
		Example: `
# Display the current configuration values
nitric config

# List all available configuration keys and their descriptions
nitric config --list
		`,
		Run: func(cmd *cobra.Command, args []string) {
			conf := do.MustInvoke[*config.Config](injector)

			underline := style.Teal("=============================")

			if listFlag {
				fmt.Println("Available configuration keys:")
				fmt.Println(underline)
				for _, key := range conf.AllKeysWithDescriptions() {
					fmt.Printf("%s: %s\n", key.Path, key.Description)
				}
				return
			}

			file := conf.FileUsed()
			if file == "" {
				file = "none (using defaults)"
			}

			fmt.Printf("Config file: %s\n", style.Teal(file))
			fmt.Println(underline)
			fmt.Println(conf.Dump())
		},
	}

	configCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all available configuration keys and their descriptions")

	var setGlobalFlag bool

	configSetCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value.",
		Example: `
# Set the URL to the Nitric server
nitric config set url https://app.nitric.io

# Set the URL to the Nitric server and save to the global config file
nitric config set url https://app.nitric.io --global
		`,
		Long: `Set a value in the cli configuration file.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]

			conf := do.MustInvoke[*config.Config](injector)

			if err := conf.SetValue(key, value); err != nil {
				fmt.Printf("error setting config: %v", err)
			}

			if err := conf.Save(setGlobalFlag); err != nil {
				fmt.Printf("error saving config: %v", err)
			}

			fmt.Printf("Config set: %s: %s\n", key, value)
		},
	}

	configSetCmd.Flags().BoolVarP(&setGlobalFlag, "global", "g", false, "Save the configuration to the global config file")

	configCmd.AddCommand(configSetCmd)
	return configCmd
}
