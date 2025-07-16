package build

import (
	"fmt"
	"slices"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/platforms"
	"github.com/nitrictech/nitric/cli/internal/plugins"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
)

type BuilderService struct {
	fs        afero.Fs
	apiClient *api.NitricApiClient
}

func (b *BuilderService) BuildProjectForTarget(appSpec *schema.Application, target string) (string, error) {
	platformRepository := platforms.NewPlatformRepository(b.apiClient)

	if len(appSpec.Targets) == 0 {
		return "", fmt.Errorf("no targets specified in project %s", appSpec.Name)
	}

	if !slices.Contains(appSpec.Targets, target) {
		return "", fmt.Errorf("target %s not found in project %s", target, appSpec.Name)
	}

	platform, err := terraform.PlatformFromId(b.fs, target, platformRepository)
	if err != nil {
		return "", err
	}

	repo := plugins.NewPluginRepository(b.apiClient)
	engine := terraform.New(platform, terraform.WithRepository(repo))

	stackPath, err := engine.Apply(appSpec)

	return stackPath, err
}

func (b *BuilderService) BuildProjectFromFileForTarget(projectFile, target string) (string, error) {
	appSpec, err := schema.LoadFromFile(b.fs, projectFile, true)
	if err != nil {
		return "", err
	}

	return b.BuildProjectForTarget(appSpec, target)
}

func NewBuilderService(injector do.Injector) (*BuilderService, error) {
	fs := do.MustInvoke[afero.Fs](injector)
	apiClient := do.MustInvoke[*api.NitricApiClient](injector)

	return &BuilderService{
		fs:        fs,
		apiClient: apiClient,
	}, nil
}
