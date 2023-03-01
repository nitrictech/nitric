package pulumi

import (
	"os/exec"
)

// TODO: Get from go.mod depedency versions
const buildkitVersion = "0.1.21"

// Installed required pulumi resources
func InstallResources() error {
	buildkitInstall := exec.Command("pulumi", "plugin", "install", "resource", "docker-buildkit", buildkitVersion, "--server", "https://github.com/MaterializeInc/pulumi-docker-buildkit/releases/download/v"+buildkitVersion)

	_, err := buildkitInstall.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
