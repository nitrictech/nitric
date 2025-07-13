package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/cli/pkg/app"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// NewTemplatesCmd creates the templates command
func NewTemplatesCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List available templates",
		Long:  `List all available templates for creating new projects.`,
		Run: func(cmd *cobra.Command, args []string) {
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.Templates())
		},
	}
}

// NewInitCmd creates the init command
func NewInitCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: fmt.Sprintf("Setup a new %s project", version.ProductName),
		Long:  fmt.Sprintf("Setup a new %s project, including within existing applications", version.ProductName),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.Init())
		},
	}
}

// NewNewCmd creates the new command
func NewNewCmd(injector do.Injector) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: fmt.Sprintf("Create a new %s project", version.ProductName),
		Long:  fmt.Sprintf("Create a new %s project from a template.", version.ProductName),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := ""
			if len(args) > 0 {
				projectName = args[0]
			}
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.New(projectName, force))
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force overwrite existing project directory")
	return cmd
}

// NewBuildCmd creates the build command
func NewBuildCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: fmt.Sprintf("Builds the %s application", version.ProductName),
		Long:  "Builds an application using the nitric.yaml application spec and referenced platform.",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.Build())
		},
	}
}

// NewGenerateCmd creates the generate command
func NewGenerateCmd(injector do.Injector) *cobra.Command {
	var (
		goFlag, pythonFlag, javascriptFlag, typescriptFlag                                    bool
		goOutputDir, goPackageName, pythonOutputDir, javascriptOutputDir, typescriptOutputDir string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: fmt.Sprintf("Generate client libraries for %s", version.ProductName),
		Long:  fmt.Sprintf("Generate client libraries for different programming languages based on the %s application specification.", version.ProductName),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.Generate(goFlag, pythonFlag, javascriptFlag, typescriptFlag, goOutputDir, goPackageName, pythonOutputDir, javascriptOutputDir, typescriptOutputDir))
		},
	}

	// Add language flags
	cmd.Flags().BoolVar(&goFlag, "go", false, "Generate Go client")
	cmd.Flags().StringVar(&goOutputDir, "go-out", "", "Output directory for Go client")
	cmd.Flags().StringVar(&goPackageName, "go-package-name", "", "Package name for Go client")

	cmd.Flags().BoolVar(&pythonFlag, "python", false, "Generate Python client")
	cmd.Flags().StringVar(&pythonOutputDir, "python-out", "", "Output directory for Python client")

	cmd.Flags().BoolVar(&javascriptFlag, "js", false, "Generate JavaScript client")
	cmd.Flags().StringVar(&javascriptOutputDir, "js-out", "", "Output directory for JavaScript client")

	cmd.Flags().BoolVar(&typescriptFlag, "ts", false, "Generate TypeScript client")
	cmd.Flags().StringVar(&typescriptOutputDir, "ts-out", "", "Output directory for TypeScript client")

	return cmd
}

// NewEditCmd creates the edit command
func NewEditCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: fmt.Sprintf("Edit the %s application", version.ProductName),
		Long:  `Edits an application using the nitric.yaml application spec and referenced platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			app, err := do.Invoke[*app.NitricApp](injector)
			if err != nil {
				cobra.CheckErr(err)
			}
			cobra.CheckErr(app.Edit())
		},
	}
}

// NewDevCmd creates the dev command
func NewDevCmd(injector do.Injector) *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: fmt.Sprintf("Run the %s application in development mode", version.ProductName),
		Long:  fmt.Sprintf("Run the %s application in development mode, allowing local testing of resources.", version.ProductName),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(app.Dev())
		},
	}
}
