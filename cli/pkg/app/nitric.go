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
	"github.com/hashicorp/go-getter"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/browser"
	"github.com/nitrictech/nitric/cli/internal/build"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/devserver"
	"github.com/nitrictech/nitric/cli/internal/simulation"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/cli/pkg/client"
	"github.com/nitrictech/nitric/cli/pkg/files"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/nitrictech/nitric/cli/pkg/tui/ask"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
)

type SugaApp struct {
	config    *config.Config
	apiClient *api.SugaApiClient
	fs        afero.Fs
	styles    tui.AppStyles
	builder   *build.BuilderService
}

func NewNitricApp(injector do.Injector) (*NitricApp, error) {
	config := do.MustInvoke[*config.Config](injector)
	apiClient := do.MustInvoke[*api.SugaApiClient](injector)
	builder := do.MustInvoke[*build.BuilderService](injector)
	fs, err := do.Invoke[afero.Fs](injector)
	if err != nil {
		fs = afero.NewOsFs()
	}

	appStyles := tui.NewAppStyles()

	return &NitricApp{config: config, apiClient: apiClient, fs: fs, builder: builder, styles: appStyles}, nil
}

// getCurrentTeam retrieves the current team from the API client
// It prints help errors to the console if the user is not authenticated or no team is set
func (c *SugaApp) getCurrentTeam() *api.Team {
	allTeams, err := c.apiClient.GetUserTeams()

	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using the", c.styles.Emphasize.Render(version.GetCommand("login")), "command")
			return nil
		}
		fmt.Printf("Failed to get teams: %v\n", err)
		return nil
	}

	var currentTeam *api.Team
	for _, t := range allTeams {
		if t.IsCurrent {
			currentTeam = &t
			break
		}
	}

	if currentTeam == nil {
		fmt.Println("No current team set, please set a team using the", c.styles.Emphasize.Render(version.GetCommand("team")), "command")
		return nil
	}

	return currentTeam
}

// Templates handles the templates command logic
func (c *SugaApp) Templates() error {
	team := c.getCurrentTeam()
	if team == nil {
		return nil
	}

	templates, err := c.apiClient.GetTemplates(team.Slug)
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using the", c.styles.Emphasize.Render(version.GetCommand("login")), "command")
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

// currentDirProjName returns a normalized version of the current directory as a name
// or "my-project" if the name can't be normalized.
func currentDirProjName() string {
	const fallback = "my-project"

	cwd, err := os.Getwd()
	if err != nil {
		return fallback
	}

	currentDir := filepath.Base(cwd)

	return normalizeDirectoryName(currentDir, fallback)
}

// Init initializes suga for an existing project, creating a suga.yaml file if it doesn't exist
func (c *SugaApp) Init() error {
	yamlPath := filepath.Join(".", version.ConfigFileName)
	exists, _ := afero.Exists(c.fs, yamlPath)

	// Read the suga.yaml file
	_, err := schema.LoadFromFile(c.fs, yamlPath, true)
	if err == nil {
		fmt.Printf("Project already initialized, run %s to edit the project\n", c.styles.Emphasize.Render(version.GetCommand("edit")))
		return nil
	} else if exists {
		fmt.Printf("Project already initialized, but an error occurred loading %s\n", c.styles.Emphasize.Render("nitric.yaml"))
		return err
	}

	fmt.Printf("Welcome to %s, this command will walk you through creating a nitric.yaml file.\n", c.styles.Emphasize.Render(version.ProductName))
	fmt.Printf("This file is used to define your app's infrastructure, resources and deployment targets.\n")
	fmt.Println()
	fmt.Printf("Here we'll only cover the basics, use %s to continue editing the project.\n", c.styles.Emphasize.Render(version.GetCommand("edit")))
	fmt.Println()

	defaultName := currentDirProjName()

	// Project Name Prompt
	var name string
	err = ask.NewInput().
		Title("Project name:").
		Value(&name).
		Placeholder(defaultName).
		Validate(func(name string) error {
			// Allow blank, we'll use the default value in this case
			if name == "" {
				return nil
			}

			return isValidProjName(name)
		}).
		Run()

	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}

	if err != nil {
		return err
	}

	if name == "" {
		name = defaultName
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

	err = schema.SaveToYaml(c.fs, yamlPath, newProject)
	if err != nil {
		fmt.Println("Failed to save " + version.ConfigFileName + " file")
		return err
	}

	fmt.Println()
	fmt.Println(c.styles.Success.Render(" " + icons.Check + " Project initialized!"))
	fmt.Println(c.styles.Faint.Render("   " + version.ProductName + " project written to " + nitricYamlPath))

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run", c.styles.Emphasize.Render(version.GetCommand("edit")), "to start the", version.ProductName, "editor")
	fmt.Println("2. Design your app's resources and deployment targets")
	fmt.Println("3. Optionally, use", c.styles.Emphasize.Render(version.GetCommand("generate")), "to generate the client libraries for your app")
	fmt.Println("4. Run", c.styles.Emphasize.Render(version.GetCommand("dev")), "to start the development server")
	fmt.Println("5. Run", c.styles.Emphasize.Render(version.GetCommand("build")), "to build the project for a specific platform")
	fmt.Println()
	fmt.Println("For more information, see the", c.styles.Emphasize.Render(version.ProductName+" docs"), "at", c.styles.Emphasize.Render("https://nitric.io/docs"))

	return nil
}

func isValidProjName(name string) error {
	// Must be kebab-case
	if !regexp.MustCompile(`^[a-z][a-z0-9-]*$`).MatchString(name) {
		return errors.New("project name must start with a letter and be lower kebab-case")
	}

	return nil
}

// normalizeDirectoryName converts a directory name to a valid project name
func normalizeDirectoryName(dirName string, fallback string) string {
	// Convert to lowercase
	normalized := strings.ToLower(dirName)

	// Replace spaces and underscores with dashes
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")

	// Remove any characters that aren't alphanumeric or dashes
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	normalized = re.ReplaceAllString(normalized, "")

	// Remove multiple consecutive dashes
	re = regexp.MustCompile(`-+`)
	normalized = re.ReplaceAllString(normalized, "-")

	// Remove leading and trailing dashes
	normalized = strings.Trim(normalized, "-")

	// If the name doesn't start with a letter, prepend "project-"
	if normalized != "" && !regexp.MustCompile(`^[a-z]`).MatchString(normalized) {
		normalized = "project-" + normalized
	}

	// If still empty or invalid, use a default
	if normalized == "" || isValidProjName(normalized) != nil {
		return fallback
	}

	return normalized
}

// New handles the new project creation command logic
func (c *SugaApp) New(projectName string, force bool) error {
	team := c.getCurrentTeam()
	if team == nil {
		return nil
	}

	templates, err := c.apiClient.GetTemplates(team.Slug)
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using", c.styles.Emphasize.Render(version.GetCommand("login")))
			return nil
		}

		return fmt.Errorf("failed to get templates: %v", err)
	}

	fmt.Printf("Welcome to %s, this command will help you create a project from a template.\n", c.styles.Emphasize.Render(version.ProductName))
	fmt.Printf("If you already have a project, run %s instead.\n", c.styles.Emphasize.Render(version.GetCommand("init")))
	fmt.Println()

	if projectName == "" {
		err := ask.NewInput().
			Title("Project name:").
			Value(&projectName).
			Validate(func(name string) error {
				if name == "" {
					return errors.New("project name is required")
				}

				return isValidProjName(name)
			}).
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

	templateDir := filepath.Join(home, version.ConfigDirName, "templates", latestVersion.TeamSlug, latestVersion.TemplateSlug, latestVersion.Version)

	templateCached, err := afero.Exists(c.fs, filepath.Join(templateDir, version.ConfigFileName))
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

	yamlPath := filepath.Join(projectDir, version.ConfigFileName)

	appSpec, err := schema.LoadFromFile(c.fs, yamlPath, false)
	if err != nil {
		return err
	}

	appSpec.Name = projectName

	err = schema.SaveToYaml(c.fs, yamlPath, appSpec)
	if err != nil {
		return err
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(c.styles.Emphasize.Render(icons.Check + " Project created!\n"))
	b.WriteString("\n")
	b.WriteString("Next steps:\n")
	b.WriteString("1. Run " + c.styles.Emphasize.Render("cd ./"+projectDir) + " to move to the project directory\n")
	b.WriteString("2. Run " + c.styles.Emphasize.Render(version.GetCommand("edit")) + " to start the " + version.ProductName + " editor\n")
	b.WriteString("3. Design your app's resources and deployment targets\n")
	b.WriteString("4. Run " + c.styles.Emphasize.Render(version.GetCommand("dev")) + " to start the development server\n")
	b.WriteString("5. Run " + c.styles.Emphasize.Render(version.GetCommand("build")) + " to build the project for a specific platform\n")
	b.WriteString("\n")
	b.WriteString("For more information, see the " + c.styles.Emphasize.Render(version.ProductName+" docs") + " at " + c.styles.Emphasize.Render("https://nitric.io/docs"))

	fmt.Println(b.String())
	return nil
}

// Build handles the build command logic
func (c *SugaApp) Build() error {
	appSpec, err := schema.LoadFromFile(c.fs, version.ConfigFileName, true)
	if err != nil {
		return err
	}

	if len(appSpec.Targets) == 0 {
		nitricEdit := c.styles.Emphasize.Render(version.GetCommand("edit"))
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

	stackPath, err := c.builder.BuildProjectForTarget(appSpec, targetPlatform)
	if err != nil {
		return err
	}

	fmt.Println(c.styles.Success.Render(" " + icons.Check + " Terraform generated successfully"))
	fmt.Println(c.styles.Faint.Render("   output written to " + stackPath))

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run", c.styles.Emphasize.Render(fmt.Sprintf("cd %s", stackPath)), "to move to the stack directory")
	fmt.Println("2. Initialize the stack", c.styles.Emphasize.Render("terraform init -upgrade"))
	fmt.Println("3. Optionally, preview with", c.styles.Emphasize.Render("terraform plan"))
	fmt.Println("4. Deploy with", c.styles.Emphasize.Render("terraform apply"))

	return nil
}

// Generate handles the generate command logic
func (c *SugaApp) Generate(goFlag, pythonFlag, javascriptFlag, typescriptFlag bool, goOutputDir, goPackageName, pythonOutputDir, javascriptOutputDir, typescriptOutputDir string) error {
	// Check if at least one language flag is provided
	if !goFlag && !pythonFlag && !javascriptFlag && !typescriptFlag {
		return fmt.Errorf("at least one language flag must be specified")
	}

	appSpec, err := schema.LoadFromFile(c.fs, version.ConfigFileName, true)
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
func (c *SugaApp) Edit() error {
	fileName := version.ConfigFileName

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

	buildServer, err := devserver.NewProjectBuild(c.apiClient, c.builder, devwsServer.Broadcast)
	if err != nil {
		return err
	}

	// create node position sync observer
	nodePositionSync := devserver.NewNodePositionSync()

	// subscribe the file sync to the websocket server
	devwsServer.Subscribe(fileSync)
	// subscribe the node position sync to the websocket server
	devwsServer.Subscribe(nodePositionSync)
	// subscribe the build server to the websocket server
	devwsServer.Subscribe(buildServer)

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	// Open browser tab to the dashboard
	devUrl := c.config.GetSugaServerUrl().JoinPath("dev")
	q := devUrl.Query()
	q.Add("port", port)
	devUrl.RawQuery = q.Encode()

	fmt.Println(tui.SugaIntro("Sync Port", port, "Dashboard", devUrl.String()))

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

	fmt.Println(c.styles.Faint.Render("Use Ctrl-C to exit"))

	// Wait for the file watcher to fail/return
	return <-errChan
}

// Dev handles the dev command logic
func Dev() error {
	// 1. Load the App Spec
	// Read the suga.yaml file
	fs := afero.NewOsFs()

	appSpec, err := schema.LoadFromFile(fs, version.ConfigFileName, true)
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
