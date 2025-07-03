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

	goOutputDir         string
	goPackageName       string
	pythonOutputDir     string
	javascriptOutputDir string
	typescriptOutputDir string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate client libraries",
	Long:  `Generate client libraries for different programming languages based on the Nitric application specification.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if at least one language flag is provided
		if !goFlag && !pythonFlag && !javascriptFlag && !typescriptFlag {
			fmt.Println("Error: at least one language flag must be specified")
			cmd.Help()
			return
		}

		fs := afero.NewOsFs()

		appSpec, err := schema.LoadFromFile(fs, "nitric.yaml", true)
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if the go language flag is provided
		if goFlag {
			fmt.Println("Generating Go client...")
			// TODO: add flags for output directory and package name
			err = client.GenerateGo(fs, *appSpec, goOutputDir, goPackageName)
			cobra.CheckErr(err)
		}

		if pythonFlag {
			fmt.Println("Generating Python client...")
			err = client.GeneratePython(fs, *appSpec, pythonOutputDir)
			cobra.CheckErr(err)
		}

		if typescriptFlag {
			fmt.Println("Generating NodeJS client...")
			err = client.GenerateTypeScript(fs, *appSpec, typescriptOutputDir)
			cobra.CheckErr(err)
		}

		fmt.Println("Clients generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Add language flags
	generateCmd.Flags().BoolVar(&goFlag, "go", false, "Generate Go client")
	generateCmd.Flags().StringVar(&goOutputDir, "go-out", "", "Output directory for Go client")
	generateCmd.Flags().StringVar(&goPackageName, "go-package-name", "", "Package name for Go client")

	generateCmd.Flags().BoolVar(&pythonFlag, "python", false, "Generate Python client")
	generateCmd.Flags().StringVar(&pythonOutputDir, "python-out", "", "Output directory for Python client")

	generateCmd.Flags().BoolVar(&javascriptFlag, "js", false, "Generate JavaScript client")
	generateCmd.Flags().StringVar(&javascriptOutputDir, "js-out", "", "Output directory for JavaScript client")

	generateCmd.Flags().BoolVar(&typescriptFlag, "ts", false, "Generate TypeScript client")
	generateCmd.Flags().StringVar(&typescriptOutputDir, "ts-out", "", "Output directory for TypeScript client")
}
