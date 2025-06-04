package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type MockTerraformPluginRepository struct {
	plugins map[string]*terraform.PluginManifest
}

func (r *MockTerraformPluginRepository) GetPlugin(name string) (*terraform.PluginManifest, error) {
	return r.plugins[name], nil
}

// TODO: remove me
func writeExampleNitricYaml(fs afero.Fs) {
	example := &schema.Application{
		Name: "test",
		ResourceIntents: map[string]schema.Resource{
			"service": {
				Type: "service",
				ServiceIntent: &schema.ServiceIntent{
					Port: 8080,
					Env: map[string]string{
						"TEST": "test",
					},
					Container: schema.Container{
						Image: &schema.DockerImage{
							ID: "test",
						},
					},
				},
			},
			"ingress": {
				Type: "entrypoint",
				EntrypointIntent: &schema.EntrypointIntent{
					Routes: map[string]schema.Route{
						"/": {
							TargetName: "service",
						},
					},
				},
			},
			"images": {
				Type:         "bucket",
				BucketIntent: &schema.BucketIntent{},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(example)
	if err != nil {
		fmt.Println("Error marshalling example to YAML:", err)
	}

	afero.WriteFile(fs, "nitric.yaml", yamlBytes, 0644)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the nitric application",
	Long:  `Builds an application using the nitric.yaml application spec and referenced platform.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Read the nitric.yaml file
		fs := afero.NewOsFs()

		// writeExampleNitricYaml(fs)

		appSpec, err := schema.LoadFromFile(fs, "nitric.yaml")
		cobra.CheckErr(err)

		// TODO: 912 repository
		embeddedRepository := terraform.NewNitricTerraformPluginRepository()

		mockPlatformRepository := terraform.NewMockPlatformRepository()

		// TODO:prompt for platform selection if multiple targets are specified
		targetPlatform := appSpec.Targets[0]

		platform, err := terraform.PlatformFromId(fs, targetPlatform, mockPlatformRepository)
		cobra.CheckErr(err)

		engine := terraform.New(platform, terraform.WithRepository(embeddedRepository))
		// Parse the application spec
		// Validate the application spec
		// Build the application using the specified platform
		// Handle any errors that occur during the build process

		err = engine.Apply(appSpec)
		if err != nil {
			fmt.Print("Error applying platform: ", err)
			return
		}

		fmt.Println("Build completed successfully.")

	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
