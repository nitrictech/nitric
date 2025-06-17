package cmd

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type MockTerraformPluginRepository struct {
	plugins map[string]*terraform.PluginManifest
}

func (r *MockTerraformPluginRepository) GetPlugin(name string) (*terraform.PluginManifest, error) {
	return r.plugins[name], nil
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the nitric application",
	Long:  `Builds an application using the nitric.yaml application spec and referenced platform.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Read the nitric.yaml file
		fs := afero.NewOsFs()

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
