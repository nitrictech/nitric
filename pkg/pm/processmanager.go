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

type process struct {
	Command []string
	cmd     *exec.Cmd
}

type pMgr struct {
	preProcesses   []*process
	userProcess    *process
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
		userProcess:    &process{Command: userCommand},
		preProcesses:   []*process{},
		monitorErrChan: make(chan error),
	}

	for _, p := range preCommands {
		m.preProcesses = append(m.preProcesses, &process{Command: p})
	}

	return m
}

func (pm *pMgr) StartUserProcess() error {
	return pm.userProcess.start()
}

func (pm *pMgr) StartPreProcesses() error {
	for i := range pm.preProcesses {
		if err := pm.preProcesses[i].start(); err != nil {
			return err
		}
	}

	return nil
}

func (pm *pMgr) StopAll() {
	err := pm.userProcess.stop()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range pm.preProcesses {
		err := p.stop()
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

func (p *process) start() error {
	if len(p.Command) == 0 {
		log.Default().Println("No Command Specified, Skipping...")

		return nil
	}

	p.cmd = exec.Command(p.Command[0], p.Command[1:]...)
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	log.Default().Printf("Starting: %s", p.Command[0])

	return errors.WithMessagef(p.cmd.Start(), "there was an error starting the process %s", p.Command[0])
}

func (p *process) stop() error {
	if p == nil || p.cmd == nil {
		return nil
	}

	err := p.cmd.Process.Signal(syscall.Signal(0))
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil
		}

		return err
	}

	err = p.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		if !errors.Is(err, os.ErrProcessDone) {
			return err
		}
	}

	return nil
}
