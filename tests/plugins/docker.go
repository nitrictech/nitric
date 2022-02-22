// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const shell = "/bin/sh"

func StartContainer(containerName string, args []string) {
	cmd := exec.Command(shell, "-c", strings.Join(args[:], " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running %s Image %v : %v \n", containerName, cmd, err)
		panic(fmt.Sprintf("Error running %s Image %v : %v", containerName, cmd, err))
	}

	// Makes process killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func StopContainer(containerName string) {
	// clean up
	stopArgs := []string{
		"docker",
		"container",
		"stop",
		containerName,
	}

	stopCmd := exec.Command(shell, "-c", strings.Join(stopArgs[:], " "))

	if err := stopCmd.Run(); err != nil {
		fmt.Printf("Error stopping %s container %v : %v \n", containerName, stopCmd, err)
		panic(fmt.Sprintf("Error stopping Firestore container %v : %v", stopCmd, err))
	}

	removeArgs := []string{
		"docker",
		"container",
		"rm",
		containerName,
	}

	removeCmd := exec.Command(shell, "-c", strings.Join(removeArgs[:], " "))

	if err := removeCmd.Run(); err != nil {
		fmt.Printf("Error removing %s container %v : %v \n", containerName, removeCmd, err)
		panic(fmt.Sprintf("Error removing Firestore container %v : %v", removeCmd, err))
	}
}
