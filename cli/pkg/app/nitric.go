package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-getter"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/browser"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/devserver"
	"github.com/nitrictech/nitric/cli/internal/platforms"
	"github.com/nitrictech/nitric/cli/internal/plugins"
	"github.com/nitrictech/nitric/cli/internal/simulation"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/nitrictech/nitric/cli/pkg/client"
	"github.com/nitrictech/nitric/cli/pkg/files"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
)

type NitricApp struct {
	config    *config.Config
	apiClient *api.NitricApiClient
	fs        afero.Fs
}

func NewNitricApp(injector do.Injector) (*NitricApp, error) {
	config := do.MustInvoke[*config.Config](injector)
	apiClient := do.MustInvoke[*api.NitricApiClient](injector)
	fs, err := do.Invoke[afero.Fs](injector)
	if err != nil {
		fs = afero.NewOsFs()
	}

	return &NitricApp{config: config, apiClient: apiClient, fs: fs}, nil
}

// Templates handles the templates command logic
func (c *NitricApp) Templates() error {
	templates, err := c.apiClient.GetTemplates()
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using the `login` command")
			fmt.Printf("%+v\n", err)
			return nil
		}

		fmt.Printf("Failed to get templates: %v\n", err)
		return nil
	}

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	fmt.Println("Available templates:")
	for _, template := range templates {
		fmt.Printf("  %s\n", template.String())
	}

	return nil
}

// Init initializes nitric for an existing project, creating a nitric.yaml file if it doesn't exist
func (c *NitricApp) Init() error {
	nitricYamlPath := filepath.Join(".", "nitric.yaml")
	exists, _ := afero.Exists(c.fs, nitricYamlPath)

	// Read the nitric.yaml file
	_, err := schema.LoadFromFile(c.fs, nitricYamlPath, true)
	if err == nil {
		fmt.Println("Project already initialized, use the `edit` command to edit the project")
		return nil
	} else if exists {
		fmt.Println("Project already initialized, but an error occurred loading nitric.yaml")
		return err
	}

	fmt.Println()
	name, err := tui.RunTextInput("Project name:", func(input string) error {
		if input == "" {
			return errors.New("project name is required")
		}

		// Must be kebab-case
		if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(input) {
			return errors.New("project name must start with a letter and be lower kebab-case")
		}

		return nil
	})
	if err != nil || name == "" {
		fmt.Println(err)
		return nil
	}

	fmt.Println()
	description, err := tui.RunTextInput("Project description (optional):", func(input string) error {
		return nil
	})

	newProject := &schema.Application{
		Name:        name,
		Description: description,
	}

	err = schema.SaveToYaml(c.fs, nitricYamlPath, newProject)
	if err != nil {
		fmt.Println("Failed to save nitric.yaml file")
		return err
	}

	successStyle := lipgloss.NewStyle().Foreground(colors.Teal).Bold(true)
	faint := lipgloss.NewStyle().Faint(true)

	fmt.Println(successStyle.Render("\nProject initialized!"))
	fmt.Println(faint.Render("nitric project written to " + nitricYamlPath))

	fmt.Println()
	_, loadEditRespIndex, err := tui.RunToggleSelect([]string{"Yes", "No"}, "Start editing in the nitric editor?")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if loadEditRespIndex == 0 {
		return c.Edit()
	}

	return nil
}

// New handles the new project creation command logic
func (c *NitricApp) New(projectName string, force bool) error {
	templates, err := c.apiClient.GetTemplates()
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using the `login` command")
			return nil
		}

		fmt.Printf("Failed to get templates: %v\n", err)
		return nil
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
			return nil
		}
	}

	projectDir := filepath.Join(".", projectName)
	if !force {
		projectExists, err := projectExists(c.fs, projectDir)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		if projectExists {
			fmt.Printf("\nDirectory ./%s already exists and is not empty\n", projectDir)
			return errors.New("project directory already exists")
		}
	}

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return errors.New("no templates available")
	}

	templateNames := make([]string, len(templates))
	for i, template := range templates {
		templateNames[i] = template.String()
	}

	// Prompt the user to select one of the templates
	fmt.Println("")
	_, index, err := tui.RunSelect(templateNames, "Template:")
	if err != nil || index == -1 {
		return err
	}

	template, err := c.apiClient.GetTemplate(templates[index].TeamSlug, templates[index].Slug, "")
	if err != nil {
		return err
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return err
	}

	templateDir := filepath.Join(home, ".nitric", "templates", template.TeamSlug, template.TemplateSlug, template.Version)

	templateCached, err := afero.Exists(c.fs, filepath.Join(templateDir, "nitric.yaml"))
	if err != nil {
		fmt.Printf("Failed read template cache directory: %v\n", err)
		return err
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
			fmt.Printf("Failed to get template: %v\n", err)
			return err
		}
	}

	// Copy the template dir contents into a new project dir
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create project directory: %v\n", err)
		return err
	}

	err = files.CopyDir(c.fs, templateDir, projectDir)
	if err != nil {
		fmt.Printf("Failed to copy template directory: %v\n", err)
		return err
	}

	nitricYamlPath := filepath.Join(projectDir, "nitric.yaml")

	appSpec, err := schema.LoadFromFile(c.fs, nitricYamlPath, false)
	if err != nil {
		return err
	}

	appSpec.Name = projectName

	err = schema.SaveToYaml(c.fs, nitricYamlPath, appSpec)
	if err != nil {
		return err
	}

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
	return nil
}

// Build handles the build command logic
func (c *NitricApp) Build() error {
	// Read the nitric.yaml file
	appSpec, err := schema.LoadFromFile(c.fs, "nitric.yaml", true)
	if err != nil {
		return err
	}

	platformRepository := platforms.NewPlatformRepository(c.apiClient)

	// TODO:prompt for platform selection if multiple targets are specified
	targetPlatform := appSpec.Targets[0]

	platform, err := terraform.PlatformFromId(c.fs, targetPlatform, platformRepository)
	if err != nil {
		return err
	}

	repo := plugins.NewPluginRepository(c.apiClient)
	engine := terraform.New(platform, terraform.WithRepository(repo))
	// Parse the application spec
	// Validate the application spec
	// Build the application using the specified platform
	// Handle any errors that occur during the build process

	err = engine.Apply(appSpec)
	if err != nil {
		fmt.Print("Error applying platform: ", err)
		return err
	}

	fmt.Println("Build completed successfully.")
	return nil
}

// Generate handles the generate command logic
func (c *NitricApp) Generate(goFlag, pythonFlag, javascriptFlag, typescriptFlag bool, goOutputDir, goPackageName, pythonOutputDir, javascriptOutputDir, typescriptOutputDir string) error {
	// Check if at least one language flag is provided
	if !goFlag && !pythonFlag && !javascriptFlag && !typescriptFlag {
		return fmt.Errorf("at least one language flag must be specified")
	}

	appSpec, err := schema.LoadFromFile(c.fs, "nitric.yaml", true)
	if err != nil {
		return err
	}

	if !client.SpecHasClientResources(*appSpec) {
		fmt.Println("No client compatible resources found in application, skipping client generation")
		return nil
	}

	// check if the go language flag is provided
	if goFlag {
		fmt.Println("Generating Go client...")
		// TODO: add flags for output directory and package name
		err = client.GenerateGo(c.fs, *appSpec, goOutputDir, goPackageName)
		if err != nil {
			return err
		}
	}

	if pythonFlag {
		fmt.Println("Generating Python client...")
		err = client.GeneratePython(c.fs, *appSpec, pythonOutputDir)
		if err != nil {
			return err
		}
	}

	if typescriptFlag {
		fmt.Println("Generating NodeJS client...")
		err = client.GenerateTypeScript(c.fs, *appSpec, typescriptOutputDir)
		if err != nil {
			return err
		}
	}

	fmt.Println("Clients generated successfully.")
	return nil
}

// Edit handles the edit command logic
func (c *NitricApp) Edit() error {
	const fileName = "nitric.yaml"

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return fmt.Errorf("error listening: %v", err)
	}

	devwsServer := devserver.NewDevWebsocketServer(devserver.WithListener(listener))
	fileSync, err := devserver.NewFileSync(fileName, devwsServer.Broadcast, devserver.WithDebounce(time.Millisecond*100))
	if err != nil {
		return err
	}
	defer fileSync.Close()

	// subscribe the file sync to the websocket server
	devwsServer.Subscribe(fileSync)

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	// Open browser tab to the dashboard
	devUrl := c.config.GetNitricServerUrl().JoinPath("dev")
	q := devUrl.Query()
	q.Add("port", port)
	devUrl.RawQuery = q.Encode()

	fmt.Println(tui.NitricIntro("Sync Port", port, "Dashboard", devUrl.String()))

	// Start the WebSocket server
	errChan := make(chan error)
	go func(errChan chan error) {
		err := devwsServer.Start()
		if err != nil {
			errChan <- err
		}
	}(errChan)

	go func() {
		err = fileSync.Start()
		if err != nil {
			fmt.Printf("Error starting file sync: %v\n", err)
		}
	}()

	fmt.Println("Opening browser to the editor")

	err = browser.Open(devUrl.String())
	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}

	// Wait for the file watcher to fail/return
	return <-errChan
}

// Dev handles the dev command logic
func Dev() error {
	// 1. Load the App Spec
	// Read the nitric.yaml file
	fs := afero.NewOsFs()

	appSpec, err := schema.LoadFromFile(fs, "nitric.yaml", true)
	if err != nil {
		return err
	}

	simserver := simulation.NewSimulationServer(fs, appSpec)
	err = simserver.Start(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

// Helper function for checking if project exists
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
