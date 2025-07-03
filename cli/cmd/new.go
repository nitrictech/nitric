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

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Nitric project",
	Run: func(cmd *cobra.Command, args []string) {

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

		// Copy the template dir contents into a new project dir
		projectDir := filepath.Join(".", projectName)
		err = os.MkdirAll(projectDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create project directory: %v", err)
			return
		}

		fs := afero.NewOsFs()

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
		b.WriteString(highlight.Render("Project created!"))
		b.WriteString("\n\n")
		b.WriteString(highlight.Render("Navigate to your project with "))
		b.WriteString("cd ./" + projectDir)
		b.WriteString("\n")
		b.WriteString("Install dependencies and you're ready to rock! 🪨")

		fmt.Println(successStyle.Render(b.String()))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
