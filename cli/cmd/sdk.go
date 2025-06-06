package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/client"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	goFlag         bool
	pythonFlag     bool
	javascriptFlag bool
	typescriptFlag bool
)

var sdkCmd = &cobra.Command{
	Use:   "sdk generate",
	Short: "Generate SDK for a specific language",
	Long:  `Generate SDK for a specific language. Supported languages: go, python, javascript, typescript.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if at least one language flag is provided
		if !goFlag && !pythonFlag && !javascriptFlag && !typescriptFlag {
			fmt.Println("Error: at least one language flag must be specified")
			cmd.Help()
			return
		}

		fs := afero.NewOsFs()

		appSpec, err := schema.LoadFromFile(fs, "nitric.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if the go language flag is provided
		if goFlag {
			fmt.Println("Generating Go SDK...")
			// TODO: add flags for output directory and package name
			err = client.GenerateGoSDK(fs, *appSpec, "", "")
			cobra.CheckErr(err)
		}

		if pythonFlag {
			fmt.Println("Generating Python SDK...")
			err = client.GeneratePythonSDK(fs, *appSpec, "", "")
			cobra.CheckErr(err)
		}

		if javascriptFlag {
			fmt.Println("Generating JavaScript SDK...")
			err = client.GenerateJavaScriptSDK(fs, *appSpec, "", "")
			cobra.CheckErr(err)
		}

		if typescriptFlag {
			fmt.Println("Generating TypeScript SDK...")
			err = client.GenerateTSSDK(fs, *appSpec, "")
			cobra.CheckErr(err)
		}

		fmt.Println("SDKs generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(sdkCmd)

	// Add language flags
	sdkCmd.Flags().BoolVar(&goFlag, "go", false, "Generate Go SDK")
	sdkCmd.Flags().BoolVar(&pythonFlag, "python", false, "Generate Python SDK")
	sdkCmd.Flags().BoolVar(&javascriptFlag, "javascript", false, "Generate JavaScript SDK")
	sdkCmd.Flags().BoolVar(&typescriptFlag, "typescript", false, "Generate TypeScript SDK")
}
