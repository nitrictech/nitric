// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/docker/docker/client"
)

type DependencyCheck func() error

func checkDockerAvailable() error {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	defer func() {
		if closeErr := cli.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	// Perform a Docker operation to verify availability
	if _, pingErr := cli.Ping(context.Background()); pingErr != nil {
		return fmt.Errorf("docker compatible API is not available, please start the docker/podman and try again")
	}

	return nil
}

func checkPulumiAvailable() error {
	_, err := exec.LookPath("pulumi")
	if err != nil {
		return fmt.Errorf("pulumi is required to use this provider, please install pulumi and try again")
	}

	return nil
}

func checkDependencies(checks ...DependencyCheck) error {
	errs := []error{}

	for _, check := range checks {
		err := check()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		errMsg := "The following dependencies are missing:"
		for _, e := range errs {
			errMsg += fmt.Sprintf("\n - %s", e.Error())
		}

		// combine the errors in a list
		return errors.New(errMsg)
	}

	return nil
}
