package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

func makeRelativeIfInCurrentDir(file string) string {
	wd, err := os.Getwd()
	if err != nil {
		return file
	}
	rel, err := filepath.Rel(wd, file)
	if err != nil {
		return file
	}
	if strings.HasPrefix(rel, "..") {
		return file
	}
	return rel
}

// NewConfigCmd creates the config command
func NewConfigCmd(injector do.Injector) *cobra.Command {
	var listFlag bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long: `Manage the CLI configuration.

By default, this command displays the current configuration values.
Use the --list flag to see all available configuration keys and their descriptions.`,
		Example: fmt.Sprintf(`
# Display the current configuration values
%s config

# List all available configuration keys and their descriptions
%s config --list
		`, version.CommandName, version.CommandName),
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
			if file != "" {
				file = makeRelativeIfInCurrentDir(file)
			} else {
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
		Example: fmt.Sprintf(`
# Set the URL to the %s server
%s config set url %s

# Set the URL to the %s server and save to the global config file
%s config set url %s --global
		`, version.ProductName, version.CommandName, version.ProductURL, version.ProductName, version.CommandName, version.ProductURL),
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
