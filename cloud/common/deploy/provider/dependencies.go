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
	"fmt"
	"os/exec"
)

type DependencyCheck func() error

func checkPulumiAvailable() error {
	_, err := exec.LookPath("pulumi")
	if err != nil {
		return fmt.Errorf("pulumi is required to use this provider, please install pulumi and try again")
	}

	return nil
}

func checkDockerAvailable() error {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("docker is required to use this provider, please install docker and try again")
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
		return fmt.Errorf(errMsg)
	}

	return nil
}
