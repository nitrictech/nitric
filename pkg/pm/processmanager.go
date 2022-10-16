// Copyright 2021 Nitric Technologies Pty Ltd.
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

package pm

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type Process struct {
	Command []string
	cmd     *exec.Cmd
}

type pMgr struct {
	preProcesses   []*Process
	userProcess    *Process
	monitorErrChan chan error
}

type ProcessManager interface {
	StartPreProcesses() error
	StartUserProcess() error
	Monitor() error
	StopAll()
}

func NewProcessManager(userCommand []string, preCommands [][]string) ProcessManager {
	m := &pMgr{
		userProcess:    &Process{Command: userCommand},
		preProcesses:   []*Process{},
		monitorErrChan: make(chan error),
	}

	for _, p := range preCommands {
		m.preProcesses = append(m.preProcesses, &Process{Command: p})
	}

	return m
}

func (pm *pMgr) StartUserProcess() error {
	if len(pm.userProcess.Command) == 0 {
		log.Default().Println("No Child Command Specified, Skipping...")

		return nil
	}

	pm.userProcess.cmd = exec.Command(pm.userProcess.Command[0], pm.userProcess.Command[1:]...)
	pm.userProcess.cmd.Stdout = os.Stdout
	pm.userProcess.cmd.Stderr = os.Stderr

	log.Default().Printf("Starting: %s", pm.userProcess.Command[0])

	return errors.WithMessagef(pm.userProcess.cmd.Start(), "there was an error starting the process %s", pm.userProcess.Command[0])
}

func (pm *pMgr) StartPreProcesses() error {
	for i := range pm.preProcesses {
		pm.preProcesses[i].cmd = exec.Command(pm.preProcesses[i].Command[0], pm.preProcesses[i].Command[1:]...)
		pm.preProcesses[i].cmd.Stdout = os.Stdout
		pm.preProcesses[i].cmd.Stderr = os.Stderr

		log.Default().Printf("Starting: %s", pm.preProcesses[i].Command[0])

		err := pm.preProcesses[i].cmd.Start()
		if err != nil {
			return errors.WithMessagef(err, "there was an error starting the process %s", pm.preProcesses[i].Command[0])
		}
	}

	return nil
}

func (pm *pMgr) stopOne(c *exec.Cmd) error {
	if c == nil {
		return nil
	}

	err := c.Process.Signal(syscall.Signal(0))
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil
		}

		return err
	}

	err = c.Process.Signal(syscall.SIGTERM)
	if err != nil {
		if !errors.Is(err, os.ErrProcessDone) {
			return err
		}
	}

	return nil
}

func (pm *pMgr) StopAll() {
	err := pm.stopOne(pm.userProcess.cmd)
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range pm.preProcesses {
		err := pm.stopOne(p.cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (pm *pMgr) Monitor() error {
	for _, p := range append(pm.preProcesses, pm.userProcess) {
		if p.cmd == nil {
			continue
		}

		go func(c *exec.Cmd) {
			pm.monitorErrChan <- c.Wait()
		}(p.cmd)
	}

	return <-pm.monitorErrChan
}
