package cmd

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/nitrictech/nitric/cli/internal/browser"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/devserver"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/spf13/cobra"
)

const fileName = "nitric.yaml"

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
		fileSync, err := devserver.NewFileSync(fileName, devwsServer.Broadcast, devserver.WithDebounce(time.Millisecond*100))
		cobra.CheckErr(err)
		defer fileSync.Close()

		// subscribe the file sync to the websocket server
		devwsServer.Subscribe(fileSync)

		port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

		// Open browser tab to the dashboard
		devUrl := config.GetNitricServerUrl().JoinPath("dev")
		q := devUrl.Query()
		q.Add("port", port)
		devUrl.RawQuery = q.Encode()

		fmt.Println(tui.NitricIntro("Sync Port", port, "Dashboard", devUrl.String()))

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

		fmt.Println("Opening browser to the editor")

		err = browser.Open(devUrl.String())
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
