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

	"github.com/charmbracelet/huh"
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
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/cli/pkg/client"
	"github.com/nitrictech/nitric/cli/pkg/files"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/nitrictech/nitric/cli/pkg/tui/ask"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
)

type NitricApp struct {
	config    *config.Config
	apiClient *api.NitricApiClient
	fs        afero.Fs
	styles    Styles
}

type Styles struct {
	emphasize lipgloss.Style
	faint     lipgloss.Style
	success   lipgloss.Style
}

func NewNitricApp(injector do.Injector) (*NitricApp, error) {
	config := do.MustInvoke[*config.Config](injector)
	apiClient := do.MustInvoke[*api.NitricApiClient](injector)
	fs, err := do.Invoke[afero.Fs](injector)
	if err != nil {
		fs = afero.NewOsFs()
	}

	return &NitricApp{config: config, apiClient: apiClient, fs: fs, styles: Styles{
		emphasize: lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
		faint:     lipgloss.NewStyle().Faint(true),
		success:   lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
	}}, nil
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
		fmt.Printf("Project already initialized, run %s to edit the project\n", c.styles.emphasize.Render(version.GetCommand("edit")))
		return nil
	} else if exists {
		fmt.Printf("Project already initialized, but an error occurred loading %s\n", c.styles.emphasize.Render("nitric.yaml"))
		return err
	}

	fmt.Printf("Welcome to %s, this command will walk you through creating a nitric.yaml file.\n", c.styles.emphasize.Render(version.ProductName))
	fmt.Printf("This file is used to define your app's infrastructure, resources and deployment targets.\n")
	fmt.Println()
	fmt.Printf("Here we'll only cover the basics, use %s to continue editing the project.\n", c.styles.emphasize.Render(version.GetCommand("edit")))
	fmt.Println()

	// Project Name Prompt
	var name string
	err = ask.NewInput().
		Title("Project name:").
		Value(&name).
		Validate(validateProjName).
		Run()

	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Project name: %s\n", name)

	// Project Description Prompt
	var description string
	err = ask.NewInput().
		Title("Project description:").
		Value(&description).
		Run()

	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Project description: %s\n", description)

	newProject := &schema.Application{
		Name:        name,
		Description: description,
	}

	err = schema.SaveToYaml(c.fs, nitricYamlPath, newProject)
	if err != nil {
		fmt.Println("Failed to save nitric.yaml file")
		return err
	}

	fmt.Println()
	fmt.Println(c.styles.success.Render(" " + icons.Check + " Project initialized!"))
	fmt.Println(c.styles.faint.Render("   " + version.ProductName + " project written to " + nitricYamlPath))

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run", c.styles.emphasize.Render(version.GetCommand("edit")), "to start the", version.ProductName, "editor")
	fmt.Println("2. Design your app's resources and deployment targets")
	fmt.Println("3. Optionally, use", c.styles.emphasize.Render(version.GetCommand("generate")), "to generate the client libraries for your app")
	fmt.Println("4. Run", c.styles.emphasize.Render(version.GetCommand("dev")), "to start the development server")
	fmt.Println("5. Run", c.styles.emphasize.Render(version.GetCommand("build")), "to build the project for a specific platform")
	fmt.Println()
	fmt.Println("For more information, see the", c.styles.emphasize.Render(version.ProductName+" docs"), "at", c.styles.emphasize.Render("https://nitric.io/docs"))

	return nil
}

func validateProjName(name string) error {
	if name == "" {
		return errors.New("project name is required")
	}

	// Must be kebab-case
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(name) {
		return errors.New("project name must start with a letter and be lower kebab-case")
	}

	return nil
}

// New handles the new project creation command logic
func (c *NitricApp) New(projectName string, force bool) error {
	templates, err := c.apiClient.GetTemplates()
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using", c.styles.emphasize.Render(version.GetCommand("login")))
			return nil
		}

		return fmt.Errorf("failed to get templates: %v", err)
	}

	fmt.Printf("Welcome to %s, this command will help you create a project from a template.\n", c.styles.emphasize.Render(version.ProductName))
	fmt.Printf("If you already have a project, run %s instead.\n", c.styles.emphasize.Render(version.GetCommand("init")))
	fmt.Println()

	if projectName == "" {
		err := ask.NewInput().
			Title("Project name:").
			Value(&projectName).
			Validate(validateProjName).
			Run()

		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}

		if err != nil {
			return err
		}
	}

	fmt.Printf("Project name: %s\n", projectName)

	projectDir := filepath.Join(".", projectName)
	if !force {
		projectExists, err := projectExists(c.fs, projectDir)
		if err != nil {
			return fmt.Errorf("failed to check if project directory exists: %v", err)
		}
		if projectExists {
			return fmt.Errorf("project directory %s already exists, use --force to overwrite", projectDir)
		}
	}

	if len(templates) == 0 {
		return errors.New("no templates found")
	}

	templateNames := make([]huh.Option[*api.Template], len(templates))
	for i, template := range templates {
		templateNames[i] = huh.NewOption(template.String(), &template)
	}

	// Prompt the user to select one of the templates
	var template *api.Template
	err = ask.NewSelect[*api.Template]().
		Title("Template:").
		Validate(func(template *api.Template) error {
			if template == nil {
				return errors.New("template is required")
			}

			return nil
		}).
		Options(templateNames...).
		Value(&template).
		Run()

	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Template: %s\n", template.String())

	latestVersion, err := c.apiClient.GetTemplate(template.TeamSlug, template.Slug, "")
	if err != nil {
		return fmt.Errorf("failed to get template: %v", err)
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	templateDir := filepath.Join(home, ".nitric", "templates", latestVersion.TeamSlug, latestVersion.TemplateSlug, latestVersion.Version)

	templateCached, err := afero.Exists(c.fs, filepath.Join(templateDir, "nitric.yaml"))
	if err != nil {
		return fmt.Errorf("failed read template cache directory: %v", err)
	}

	if !templateCached {
		goGetter := &getter.Client{
			Ctx:             context.Background(),
			Dst:             templateDir,
			Src:             latestVersion.GitSource,
			Mode:            getter.ClientModeAny,
			DisableSymlinks: true,
		}

		err = goGetter.Get()
		if err != nil {
			return fmt.Errorf("failed to get template: %v", err)
		}
	}

	// Copy the template dir contents into a new project dir
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	err = files.CopyDir(c.fs, templateDir, projectDir)
	if err != nil {
		return fmt.Errorf("failed to copy template directory: %v", err)
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

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(c.styles.emphasize.Render(icons.Check + " Project created!\n"))
	b.WriteString("\n")
	b.WriteString("Next steps:\n")
	b.WriteString("1. Run " + c.styles.emphasize.Render("cd ./"+projectDir) + " to move to the project directory\n")
	b.WriteString("2. Run " + c.styles.emphasize.Render(version.GetCommand("edit")) + " to start the " + version.ProductName + " editor\n")
	b.WriteString("3. Design your app's resources and deployment targets\n")
	b.WriteString("4. Run " + c.styles.emphasize.Render(version.GetCommand("dev")) + " to start the development server\n")
	b.WriteString("5. Run " + c.styles.emphasize.Render(version.GetCommand("build")) + " to build the project for a specific platform\n")
	b.WriteString("\n")
	b.WriteString("For more information, see the " + c.styles.emphasize.Render(version.ProductName+" docs") + " at " + c.styles.emphasize.Render("https://nitric.io/docs"))

	fmt.Println(b.String())
	return nil
}

// Build handles the build command logic
func (c *NitricApp) Build() error {
	appSpec, err := schema.LoadFromFile(c.fs, "nitric.yaml", true)
	if err != nil {
		return err
	}

	platformRepository := platforms.NewPlatformRepository(c.apiClient)

	if len(appSpec.Targets) == 0 {
		nitricEdit := c.styles.emphasize.Render(version.GetCommand("edit"))
		fmt.Printf("No targets specified in nitric.yaml, run %s to add a target\n", nitricEdit)
		return nil
	}

	var targetPlatform string

	if len(appSpec.Targets) == 1 {
		targetPlatform = appSpec.Targets[0]
	} else {
		err := ask.NewSelect[string]().
			Title("Select a build target").
			Options(huh.NewOptions(appSpec.Targets...)...).
			Value(&targetPlatform).
			Validate(func(targetPlatform string) error {
				if targetPlatform == "" {
					return errors.New("target platform is required")
				}

				return nil
			}).
			Run()

		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}

		if err != nil {
			return err
		}
	}

	if targetPlatform == "" {
		return fmt.Errorf("no target platform selected")
	}

	platform, err := terraform.PlatformFromId(c.fs, targetPlatform, platformRepository)
	if errors.Is(err, terraform.ErrUnauthenticated) {
		fmt.Printf("Please login first, using the %s command\n", c.styles.emphasize.Render(version.GetCommand("login")))
		return nil
	} else if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\nBuilding application for %s\n\n", c.styles.emphasize.Render(targetPlatform))

	repo := plugins.NewPluginRepository(c.apiClient)
	engine := terraform.New(platform, terraform.WithRepository(repo))

	stackPath, err := engine.Apply(appSpec)
	if err != nil {
		fmt.Print("Error applying platform: ", err)
		return err
	}

	fmt.Println(c.styles.success.Render(" " + icons.Check + " Terraform generated successfully"))
	fmt.Println(c.styles.faint.Render("   output written to " + stackPath))

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run", c.styles.emphasize.Render(fmt.Sprintf("cd %s", stackPath)), "to move to the stack directory")
	fmt.Println("2. Initialize the stack", c.styles.emphasize.Render("terraform init -upgrade"))
	fmt.Println("3. Optionally, preview with", c.styles.emphasize.Render("terraform plan"))
	fmt.Println("4. Deploy with", c.styles.emphasize.Render("terraform apply"))

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
