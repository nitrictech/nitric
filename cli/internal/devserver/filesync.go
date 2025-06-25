package devserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type NitricFileSync struct {
	filePath         string
	file             *os.File
	debounce         time.Duration
	publish          PublishFunc
	lastSyncContents []byte
}

type FileSyncMessage Message[schema.Application]

type FileSyncOption func(*NitricFileSync)

func WithDebounce(debounce time.Duration) FileSyncOption {
	return func(fs *NitricFileSync) {
		fs.debounce = debounce
	}
}

func (fs *NitricFileSync) OnMessage(message json.RawMessage) {
	var fileSyncMessage FileSyncMessage
	err := json.Unmarshal(message, &fileSyncMessage)
	if err != nil {
		return
	}

	// Not the right message type continue
	if fileSyncMessage.Type != "nitricSync" {
		return
	}

	yamlContents, err := yaml.Marshal(fileSyncMessage.Payload)
	if err != nil {
		fmt.Println("Error marshalling application to yaml:", err)
		return
	}

	_, err = fs.file.Seek(0, 0)
	if err != nil {
		return
	}
	err = fs.file.Truncate(0)
	if err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}

	fs.lastSyncContents = yamlContents

	_, err = fs.file.Write(yamlContents)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func (fs *NitricFileSync) Close() error {
	return fs.file.Close()
}

func NewFileSync(filePath string, publish PublishFunc, options ...FileSyncOption) (*NitricFileSync, error) {
	var err error
	fs := &NitricFileSync{
		filePath: filePath,
		publish:  publish,
	}

	for _, option := range options {
		option(fs)
	}

	fs.file, err = os.OpenFile(fs.filePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *NitricFileSync) Start() error {
	return fs.watchFile()
}

// watchFile watches the file for changes and broadcasts updates
func (fw *NitricFileSync) watchFile() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Add the file to the watcher
	err = watcher.Add(fw.filePath)
	if err != nil {
		return err
	}

	var cancel, debounced func()
	for event := range watcher.Events {
		if event.Has(fsnotify.Write) {
			if cancel != nil {
				cancel()
			}

			var fileError error = nil
			debounced, cancel = lo.NewDebounce(fw.debounce, func() {
				// Read the current contents of the file
				fw.file.Seek(0, 0) // Seek to beginning
				contents, err := io.ReadAll(fw.file)
				if err != nil {
					fileError = err
					return
				}

				application, schemaResult, err := schema.ApplicationFromYaml(string(contents))
				if err != nil {
					fmt.Println("Error parsing application from yaml:", err)
					return
				} else if schemaResult != nil && len(schemaResult.Errors()) > 0 {
					fmt.Println("Errors parsing application from yaml:", schemaResult.Errors())
					return
				}

				// Don't publish if the contents are the same as the last sync down
				if bytes.Equal(fw.lastSyncContents, contents) {
					return
				}

				fw.publish(Message[any]{
					Type:    "nitricSync",
					Payload: *application,
				})
			})
			debounced()
			if fileError != nil {
				return err
			}
		}
	}

	return nil
}
