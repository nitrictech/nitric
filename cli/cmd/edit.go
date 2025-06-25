package cmd

import (
	"log"
	"time"

	"github.com/nitrictech/nitric/cli/internal/filex"
	"github.com/spf13/cobra"
)

const fileName = "nitric.yaml"

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the nitric application",
	Long:  `Edits an application using the nitric.yaml application spec and referenced platform.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create and start the file watcher
		fileWatcher := filex.NewWebsocketServerSync(fileName, filex.WithDebounce(time.Millisecond*100))
		defer fileWatcher.Close()

		log.Printf("Watching file: %s", fileName)
		log.Println("Press Ctrl+C to stop")

		// Start the WebSocket server
		err := fileWatcher.Start()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
