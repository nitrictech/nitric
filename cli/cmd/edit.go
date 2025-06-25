package cmd

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/browser"
	"github.com/nitrictech/nitric/cli/internal/devserver"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/spf13/cobra"
)

const fileName = "nitric.yaml"

var nitricSimple = style.Purple("Nitric")

var nitricLightning = style.Purple(icons.Lightning + " Nitric")

func nitricIntro() string {
	version := version.GetShortVersion()

	intro := fmt.Sprintf("\n%s %s\n", nitricLightning, style.Gray(version))

	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder(), false, true).Render(intro)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the nitric application",
	Long:  `Edits an application using the nitric.yaml application spec and referenced platform.`,
	Run: func(cmd *cobra.Command, args []string) {

		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			log.Printf("Error listening: %v", err)
		}

		devwsServer := devserver.NewDevWebsocketServer(devserver.WithListener(listener))
		fileSync, err := devserver.NewFileSync(fileName, devwsServer.Publish, devserver.WithDebounce(time.Millisecond*100))
		cobra.CheckErr(err)
		defer fileSync.Close()

		// subscribe the file sync to the websocket server
		devwsServer.Subscribe(fileSync)

		fmt.Println(nitricIntro())

		fmt.Println("Starting nitric.yaml synchronizer on port", listener.Addr().(*net.TCPAddr).Port)

		// Start the WebSocket server
		errChan := make(chan error)
		go func(errChan chan error) {
			err := devwsServer.Start()
			if err != nil {
				errChan <- err
			}
		}(errChan)

		go func() {
			err = fileSync.Start()
			if err != nil {
				log.Printf("Error starting file sync: %v", err)
			}
		}()

		// Get the port for the listener
		port := listener.Addr().(*net.TCPAddr).Port

		fmt.Printf("Opening browser tab to the %s editor\n", nitricSimple)

		// Open browser tab to the dashboard
		err = browser.Open(fmt.Sprintf("http://localhost:8080/dev?port=%d", port))
		if err != nil {
			log.Printf("Error opening browser: %v", err)
		}

		// Wait for the file watcher to fail/return
		cobra.CheckErr(<-errChan)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
