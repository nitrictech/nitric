package cmd

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nitrictech/nitric/cli/internal/browser"
	"github.com/nitrictech/nitric/cli/internal/filex"
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

		// Create and start the file watcher
		fileWatcher := filex.NewWebsocketServerSync(fileName, filex.WithDebounce(time.Millisecond*100), filex.WithListener(listener))
		defer fileWatcher.Close()

		log.Printf("Watching file: %s", fileName)
		log.Println("Press Ctrl+C to stop")

		// Start the WebSocket server
		errChan := make(chan error)
		go func(errChan chan error) {
			err := fileWatcher.Start()
			if err != nil {
				errChan <- err
			}
		}(errChan)

		// Get the port for the listener
		port := listener.Addr().(*net.TCPAddr).Port

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
