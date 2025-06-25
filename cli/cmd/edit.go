package cmd

import (
	"log"

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
		fileWatcher := filex.NewFileSyncer(fileName)
		defer fileWatcher.Close()

		// Start the WebSocket server
		err := fileWatcher.StartServer()
		cobra.CheckErr(err)

		log.Printf("Watching file: %s", fileName)
		log.Println("Press Ctrl+C to stop")

		// Block forever
		select {}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
