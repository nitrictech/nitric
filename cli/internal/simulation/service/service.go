package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/nitrictech/nitric/cli/internal/netx"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/robfig/cron/v3"
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

	port    netx.ReservedPort
	apiPort netx.ReservedPort

	// Rapid failure detection
	consecutiveFailures int // Count of consecutive rapid failures
	maxFailures         int // Maximum consecutive failures allowed
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

// hasExceededFailureLimit checks if we've had too many consecutive failures
func (s *ServiceSimulation) hasExceededFailureLimit() bool {
	if s.maxFailures <= 0 {
		return false // No limit configured
	}
	return s.consecutiveFailures >= s.maxFailures
}

func (s *ServiceSimulation) startSchedules(stdoutWriter, stderrorWriter io.Writer) (*cron.Cron, error) {
	triggers := s.intent.Triggers
	cron := cron.New()

	for triggerName, trigger := range triggers {
		if trigger.Schedule != nil {
			url := url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("localhost:%d", s.port),
				Path:   trigger.Path,
			}

			_, err := cron.AddFunc(trigger.Schedule.CronExpression, func() {
				req, err := http.NewRequest(http.MethodPost, url.String(), nil)
				if err != nil {
					// log the error
					fmt.Fprint(stderrorWriter, "error creating request for schedule", err)
					return
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprint(stderrorWriter, "error sending request for schedule", err)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Fprint(stderrorWriter, "error sending request for schedule", resp.StatusCode)
					return
				}

				fmt.Fprintf(stdoutWriter, "schedule [%s] triggered on %s", triggerName, trigger.Path)
			})

			if err != nil {
				return nil, err
			}
		}
	}

	cron.Start()

	return cron, nil
}

func (s *ServiceSimulation) Start(autoRestart bool) error {
	s.autoRestart = autoRestart

	stdoutWriter := newServiceLogWriter(s, OutputType_Stdout)
	stderrWriter := newServiceLogWriter(s, OutputType_Stderr)

	// Track if service ran successfully for a meaningful duration
	const minSuccessfulRuntime = 5 * time.Second

	for {
		if s.currentStatus != Status_Init && !s.autoRestart {
			break
		}

		// Check if we've exceeded the failure limit before attempting to restart
		if s.hasExceededFailureLimit() {
			fmt.Fprintf(stderrWriter, "\n[ERROR] Service '%s' has failed %d consecutive times. Stopping restart attempts to prevent crash loop.\n",
				s.name, s.maxFailures)
			s.updateStatus(Status_Fatal)
			return fmt.Errorf("service '%s' exceeded maximum consecutive restart attempts (%d failures)",
				s.name, s.maxFailures)
		}

		commandParts := strings.Split(s.intent.Dev.Script, " ")

		// The $PORT env is also needed if the dev command decides the port the service should listen on.
		// e.g. python `uv run uvicorn main:app --host 0.0.0.0 --port $PORT`
		for i := range commandParts {
			commandParts[i] = strings.ReplaceAll(commandParts[i], "$PORT", fmt.Sprintf("%d", s.port))
		}

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
		srvCommand.Env = append(srvCommand.Env, fmt.Sprintf("SUGA_SERVICE_ADDRESS=localhost:%d", s.apiPort))

		srvCommand.Dir = s.intent.Container.Docker.Context
		srvCommand.Stdout = stdoutWriter
		srvCommand.Stderr = stderrWriter

		startTime := time.Now()
		err := srvCommand.Start()
		if err != nil {
			// If the command fails to start at all (e.g., command not found, compilation error),
			// this is an immediate fatal error - no point retrying
			fmt.Fprintf(stderrWriter, "\n Service failed to start: %v\n", err)
			s.updateStatus(Status_Fatal)
			return fmt.Errorf("service %s failed to start: %w", s.name, err)
		}

		s.cmd = srvCommand
		s.updateStatus(Status_Running)

		cron, err := s.startSchedules(stdoutWriter, stderrWriter)
		if err != nil {
			s.updateStatus(Status_Fatal)
			return err
		}

		err = srvCommand.Wait()
		runDuration := time.Since(startTime)

		// Stop running cron jobs
		cron.Stop()

		if err == nil {
			// Clean exit
			break
		}

		// Check if the service ran for a meaningful duration
		if runDuration < minSuccessfulRuntime {
			// Service crashed quickly - likely a startup issue
			s.consecutiveFailures++
			fmt.Fprintf(stderrWriter, "\nService crashed after only %v\n", runDuration)
		} else {
			// Service ran for a while before crashing - reset failure count
			s.consecutiveFailures = 0
		}

		s.updateStatus(Status_Stopped)
	}

	return nil
}

func NewServiceSimulation(name string, intent schema.ServiceIntent, port netx.ReservedPort, apiPort netx.ReservedPort) (*ServiceSimulation, <-chan ServiceEvent, error) {
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
		apiPort:       apiPort,
		maxFailures:   2, // Default: stop after 1 consecutive rapid failures
	}, eventsChan, nil
}
