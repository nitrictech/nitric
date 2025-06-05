package cmd

import (
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/cli/pkg/sdk"
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

		langs := []sdk.Language{}
		if goFlag {
			langs = append(langs, sdk.Go)
		}
		if pythonFlag {
			langs = append(langs, sdk.Python)
		}
		if javascriptFlag {
			langs = append(langs, sdk.Javascript)
		}
		if typescriptFlag {
			langs = append(langs, sdk.Typescript)
		}

		langStrings := make([]string, len(langs))
		for i, lang := range langs {
			langStr, err := sdk.LanguageToString(lang)
			if err != nil {
				fmt.Println("Error converting language to string:", err)
				return
			}
			langStrings[i] = langStr
		}

		fmt.Println("Generating SDKs for languages:", strings.Join(langStrings, ", "))

		err = sdk.GenerateSDKs(fs, *appSpec, "", langs)
		if err != nil {
			fmt.Println("Error generating SDKs:", err)
			return
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
