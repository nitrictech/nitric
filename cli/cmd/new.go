package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-getter"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/auth"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/nitrictech/nitric/cli/pkg/files"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func projectExists(fs afero.Fs, projectDir string) (bool, error) {
	projectExists, err := afero.Exists(fs, projectDir)
	if err != nil {
		return false, fmt.Errorf("failed to read intended project directory: %v", err)
	}
	if projectExists {
		// Check if the directory is empty
		files, err := afero.ReadDir(fs, projectDir)
		if err != nil {
			return false, fmt.Errorf("failed to read project directory: %v", err)
		}
		return len(files) > 0, nil
	}
	return false, nil
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Nitric project",
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Printf("Failed to get force flag: %v", err)
			return
		}

		fs := afero.NewOsFs()

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		if projectName == "" {
			fmt.Println()
			var err error
			projectName, err = tui.RunTextInput("Project name:", func(input string) error {
				if input == "" {
					return errors.New("project name is required")
				}

				// Must be kebab-case
				if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(input) {
					return errors.New("project name must start with a letter and be lower kebab-case")
				}

				return nil
			})
			if err != nil || projectName == "" {
				fmt.Println(err)
				fmt.Println("+" + projectName + "+")
				return
			}
		}

		projectDir := filepath.Join(".", projectName)
		if !force {
			projectExists, err := projectExists(fs, projectDir)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if projectExists {
				fmt.Printf("\nDirectory ./%s already exists and is not empty\n", projectDir)
				return
			}
		}

		token, err := auth.GetOrRefreshWorkosToken()
		if err != nil {
			fmt.Printf("\n Not currently logged in, run `nitric login` to login")
			return
		}

		client := api.NewNitricApiClient(config.GetNitricServerUrl(), &token.AccessToken)

		resp, err := client.GetTemplates()
		if err != nil {
			fmt.Printf("Failed to get templates: %v", err)
			return
		}

		if len(resp) == 0 {
			fmt.Println("No templates found")
			return
		}

		templateNames := make([]string, len(resp))
		for i, template := range resp {
			templateNames[i] = template.String()
		}

		// Prompt the user to select one of the templates
		fmt.Println("")
		_, index, err := tui.RunSelect(templateNames, "Template:")
		if err != nil || index == -1 {
			return
		}

		template, err := client.GetTemplate(resp[index].TeamSlug, resp[index].Slug, "")
		cobra.CheckErr(err)

		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			return
		}

		templateDir := filepath.Join(home, ".nitric", "templates", template.TeamSlug, template.TemplateSlug, template.Version)

		templateCached, err := afero.Exists(fs, filepath.Join(templateDir, "nitric.yaml"))
		if err != nil {
			fmt.Printf("Failed read template cache directory: %v", err)
			return
		}

		if !templateCached {
			goGetter := &getter.Client{
				Ctx:             context.Background(),
				Dst:             templateDir,
				Src:             template.GitSource,
				Mode:            getter.ClientModeAny,
				DisableSymlinks: true,
			}

			err = goGetter.Get()
			if err != nil {
				fmt.Printf("Failed to get template: %v", err)
				return
			}
		}

		// Copy the template dir contents into a new project dir
		err = os.MkdirAll(projectDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create project directory: %v", err)
			return
		}

		err = files.CopyDir(fs, templateDir, projectDir)
		if err != nil {
			fmt.Printf("Failed to copy template directory: %v", err)
			return
		}

		nitricYamlPath := filepath.Join(projectDir, "nitric.yaml")

		appSpec, err := schema.LoadFromFile(fs, nitricYamlPath, false)
		cobra.CheckErr(err)

		appSpec.Name = projectName

		err = schema.SaveToYaml(fs, nitricYamlPath, appSpec)
		cobra.CheckErr(err)

		successStyle := lipgloss.NewStyle().MarginLeft(3)
		highlight := lipgloss.NewStyle().Foreground(colors.Teal).Bold(true)

		var b strings.Builder

		b.WriteString("\n")
		b.WriteString("Project created!")
		b.WriteString("\n\n")
		b.WriteString("Navigate to your project with ")
		b.WriteString(highlight.Render("cd ./" + projectDir))
		b.WriteString("\n")
		b.WriteString("Install dependencies and you're ready to rock! 🪨")

		fmt.Println(successStyle.Render(b.String()))
	},
}

func init() {
	newCmd.Flags().BoolP("force", "f", false, "Force overwrite existing project directory")
	rootCmd.AddCommand(newCmd)
}
