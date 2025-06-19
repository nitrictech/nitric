package service

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"github.com/nitrictech/nitric/cli/internal/netx"
	"github.com/nitrictech/nitric/cli/pkg/schema"
)

type ServiceSimulation struct {
	name   string
	intent schema.ServiceIntent
	cmd    *exec.Cmd

	currentStatus Status

	autoRestart bool

	events chan<- ServiceEvent

	stdOut io.Writer
	stdErr io.Writer

	port netx.ReservedPort
}

var _ SimulatedService = (*ServiceSimulation)(nil)

type ServiceEvent struct {
	SimulatedService
	PreviousStatus Status
	Output         *OutputType
	Content        []byte
}

type OutputType bool

const (
	OutputType_Stdout OutputType = false
	OutputType_Stderr OutputType = true
)

type SimulatedService interface {
	GetName() string
	GetPort() int
	GetCmd() *exec.Cmd
	GetStatus() Status
	Signal(sig os.Signal)
}

func (s *ServiceSimulation) GetCmd() *exec.Cmd {
	return s.cmd
}

func (s *ServiceSimulation) GetPort() int {
	return int(s.port)
}

func (s *ServiceSimulation) GetName() string {
	return s.name
}

var stopSignals = []os.Signal{
	syscall.SIGABRT,
	syscall.SIGALRM,
	syscall.SIGTERM,
	os.Interrupt,
}

// Signal - sends a signal to the service process
func (s *ServiceSimulation) Signal(sig os.Signal) {
	if slices.Contains(stopSignals, sig) {
		s.autoRestart = false
		s.updateStatus(Status_Stopping)
	}
	// If windows, it will always Kill 🔪... (signals are not supported on windows)
	err := s.cmd.Process.Signal(sig)
	if err != nil {
		s.autoRestart = false
		s.updateStatus(Status_Stopping)
		err = s.cmd.Process.Kill()
	}
}

func (s *ServiceSimulation) GetStatus() Status {
	return s.currentStatus
}

func (s *ServiceSimulation) PublishEvent(event ServiceEvent) {
	s.events <- event
}

func (s *ServiceSimulation) updateStatus(newStatus Status) {
	previousStatus := s.currentStatus
	s.currentStatus = newStatus

	s.PublishEvent(ServiceEvent{
		SimulatedService: s,
		PreviousStatus:   previousStatus,
	})
}

func (s *ServiceSimulation) Start(autoRestart bool) error {
	s.autoRestart = autoRestart

	stdoutWriter := newServiceLogWriter(s, OutputType_Stdout)
	stderrWriter := newServiceLogWriter(s, OutputType_Stderr)

	for {
		if s.currentStatus != Status_Init && !s.autoRestart {
			break
		}

		commandParts := strings.Split(s.intent.Dev.Script, " ")
		srvCommand := exec.Command(
			commandParts[0],
			commandParts[1:]...,
		)

		srvCommand.Env = append([]string{}, os.Environ()...)

		if s.currentStatus == Status_Init {
			s.updateStatus(Status_Starting)
		} else {
			s.updateStatus(Status_Restarting)
		}

		srvCommand.Env = append(srvCommand.Env, fmt.Sprintf("PORT=%d", s.port))

		srvCommand.Dir = s.intent.Container.Docker.Context
		srvCommand.Stdout = stdoutWriter
		srvCommand.Stderr = stderrWriter

		err := srvCommand.Start()
		if err != nil {
			s.updateStatus(Status_Fatal)
			return err
		}
		s.updateStatus(Status_Running)

		err = srvCommand.Wait()
		if err == nil {
			break
		}

		s.updateStatus(Status_Stopped)
	}

	return nil
}

func NewServiceSimulation(name string, intent schema.ServiceIntent, port netx.ReservedPort) (*ServiceSimulation, <-chan ServiceEvent, error) {
	if intent.Dev == nil {
		return nil, nil, fmt.Errorf("service does not have a dev configuration and cannot be started")
	}

	if intent.Dev.Script == "" {
		return nil, nil, fmt.Errorf("service does not have a dev script and cannot be started")
	}

	eventsChan := make(chan ServiceEvent)

	return &ServiceSimulation{
		name:          name,
		intent:        intent,
		currentStatus: Status_Init,
		events:        eventsChan,
		port:          port,
	}, eventsChan, nil
}
