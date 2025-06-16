package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/netx"
	"github.com/nitrictech/nitric/cli/internal/simulation"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func (p *PrefixWriter) Write(content []byte) (int, error) {
	_, err := fmt.Fprintf(p.writer, "%s: %s", p.prefix, string(content))
	if err != nil {
		return 0, err
	}

	return len(content), nil
}

func NewPrefixWriter(prefix string, writer io.Writer) *PrefixWriter {
	return &PrefixWriter{
		prefix: prefix,
		writer: writer,
	}
}

var someChars = "❍ ⚬ ♽ ♼ ☃ ⚙  ⚞ ⚟ ⚠"

var (
	rightArrow   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(14)).Render("➜")
	recycle      = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(14)).Render("♻")
	warning      = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(14)).Render("⚠")
	svcStyle     = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	styledNitric = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(13)).Render("nitric")
)

var dev = &cobra.Command{
	Use:   "dev",
	Short: "Run the Nitric application in development mode",
	Long:  `Run the Nitric application in development mode, allowing for live reloading and local testing of resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Extract common loading logic into a separate function
		// (see build.go for common loading logic)

		// 1. Load the App Spec
		// Read the nitric.yaml file
		fs := afero.NewOsFs()

		appSpec, err := schema.LoadFromFile(fs, "nitric.yaml")
		cobra.CheckErr(err)

		// 3. Run additional services (docker compose?)
		waitGroup, _ := errgroup.WithContext(context.Background())

		waitGroup.Go(func() error {
			// 4. Start default local simulation
			simserver := simulation.NewSimulationServer(fs)
			return simserver.Start()
		})

		var done atomic.Bool
		done.Store(false)

		// 5. Start services
		runningProcessesLock := sync.Mutex{}
		runningProcesses := map[string]*exec.Cmd{}
		services := appSpec.GetServiceIntents()
		for serviceName, intent := range services {
			styledSvcName := svcStyle.Render(serviceName)

			logwriter := NewPrefixWriter(styledSvcName, os.Stdout)

			waitGroup.Go(func() error {
				// Start the service command, restarting if it closes/crashes
				for {
					if intent.Dev.Script == "" {
						break // Skip services without a dev script
					}

					commandParts := strings.Split(intent.Dev.Script, " ")
					srvCommand := exec.Command(
						commandParts[0],
						commandParts[1:]...,
					)

					srvCommand.Env = append([]string{}, os.Environ()...)

					// Get an available port for the service
					port, err := netx.GetNextPort()
					cobra.CheckErr(err)

					fmt.Printf("%s Starting service %s on port %d\n", rightArrow, styledSvcName, port)

					srvCommand.Env = append(srvCommand.Env, fmt.Sprintf("PORT=%d", port))

					srvCommand.Dir = intent.Container.Docker.Context
					srvCommand.Stdout = logwriter

					err = srvCommand.Start()
					if err != nil {
						return err
					}

					runningProcessesLock.Lock()
					runningProcesses[serviceName] = srvCommand
					runningProcessesLock.Unlock()

					err = srvCommand.Wait()
					if err == nil || done.Load() {
						break
					}

					runningProcessesLock.Lock()
					delete(runningProcesses, serviceName)
					runningProcessesLock.Unlock()
					// try to restart the process
					fmt.Printf("Service %s: %v\n", styledSvcName, err)
				}

				return nil
			})

		}

		errChan := make(chan error)
		go func() {
			err := waitGroup.Wait()
			errChan <- err
		}()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

		select {
		case err := <-errChan:
			// Never try to restart services in this case
			done.Store(true)
			if err != nil {
				log.Println(err)
			}
		case sig := <-sigChan:
			// Never try to restart services in this case
			done.Store(true)
			for _, srvCmd := range runningProcesses {
				// If windows, it will always Kill 🔪... (signals are not supported on windows)
				err := srvCmd.Process.Signal(sig)
				if err != nil {
					err = srvCmd.Process.Kill()
				}
			}
		}
	},
}

func init() {
	// Add the dev command to the root command
	rootCmd.AddCommand(dev)

	// Add flags for the dev command if needed
	// e.g., dev.Flags().StringVarP(&flagName, "flag", "f", "defaultValue", "Description of the flag")
}

//
